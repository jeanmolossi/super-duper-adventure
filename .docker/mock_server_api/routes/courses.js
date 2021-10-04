const { Router } = require("express");
const faker = require("faker");
const { courses, createCurso } = require("../services/courses");
const { students } = require("../services/students");
const { rabbitMQ } = require("../services/rabbit");
const { sleep } = require("../utils/sleep");

const coursesResource = Router();

/** POST Create course */
coursesResource.post('/courses', (req, res) => {
	const { id } = req.body;
	const course = courses.create(createCurso(id.toString()));
	return res.json(course)
})

/** GET Get courses */
coursesResource.get('/courses', (req, res) => {
	const { query } = req;
	const ids = query.ids?.split(',')
	const coursesResponse = courses.index(ids);

	if (!coursesResponse)
		return res.status(404).json({ message: "Courses not found" })

	return res.json({
		total: coursesResponse.length,
		data: coursesResponse
	})
})

/** GET Get course by ID */
coursesResource.get('/course/:id', async (req, res) => {
	const { id } = req.params;
	const course = courses.show(id.toString())

	if (!course) {
		return res.status(404).json({ message: "Course not found" })
	}

	if (!req.headers['Origin-App'] && !req.headers['origin-app']) {
		console.log('PUBLISH')
		await rabbitMQ.publish(course.id.toString())
	}

	return res.json(course)
})

/** GET Get students */
coursesResource.get('/course/:id/students', async (req, res) => {
	const { id } = req.params;

	const course = courses.show(id.toString())
	const courseStudents = students.indexFromCourse(id.toString())

	const data = {
		...course,
		alunos: courseStudents
	}

	const fakeLatency = faker.datatype.number({ min: 25, max: 1200 });
	await sleep(fakeLatency);

	return res.json(data)
})

module.exports = coursesResource;