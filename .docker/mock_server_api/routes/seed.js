const { Router } = require("express");
const faker = require("faker");
const _ = require("lodash");
const { students, createUser } = require("../services/students");
const { courses, createCurso } = require("../services/courses");
const { rabbitMQ } = require("../services/rabbit");

const seedResource = Router();

function extractTotal(req, max = 65000, rand = false) {
	const { query } = req;

	let total = query.total;

	if (!total) {
		total = faker.datatype.number({ min: 50, max: 250 })
	}

	if (total > max) {
		total = max;
	}

	if (rand) {
		const min = +total > 5000 ? +total - 5000 : +total;
		const max = +total + 25000;
		total = faker.datatype.number({ min, max })
	}

	return total;
}

seedResource.get('/students', (req, res) => {
	const total = extractTotal(req);

	const coursesResponse = _.shuffle(courses.index());

	let counter = 0;
	coursesResponse.forEach(course => {
		for (let i = 0; i < total; i++) {
			counter++;
			students.create(createUser(course.id))
		}
	})

	return res.json({
		message: `seeded-students. total ${ counter } records`,
		total: counter
	})
})

seedResource.get('/courses', (req, res) => {
	const total = extractTotal(req, 5000, true);

	for (let i = 0; i < total; i++) {
		courses.create(createCurso(+i))
	}

	return res.json({
		message: `seeded-courses. total ${ total } records`,
		total
	})
})
seedResource.get('/queue', async (req, res) => {
	const total = extractTotal(req, 65500, true);

	const fetchCourses = _.shuffle(courses.index());

	const coursesIds = fetchCourses.length <= total
		? fetchCourses.map(c => c.id.toString())
		: Array.from({ length: total }).map((_, i) => fetchCourses[i].id.toString())

	const promises = coursesIds.map(async (id) => rabbitMQ.publish(id))

	try {
		await Promise.all(promises)
	} catch (e) {
		console.log("Falhou ao publicar ", e.message)
	}

	const messagesPublished = coursesIds.length;

	const data = {
		message: `seeded-queue. Total ${ messagesPublished } records`,
		total: messagesPublished
	}

	return res.json(data)

})

module.exports = seedResource