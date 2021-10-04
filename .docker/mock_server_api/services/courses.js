const faker = require("faker");
const _ = require("lodash");

class Courses {

	#courses = []

	create(payload) {
		this.#courses.push(payload)

		return payload;
	}

	show(id) {
		if (!id) return;

		const course = _.find(this.#courses, { id });

		return course || {
			id,
			code: 404,
			message: "Course not found"
		}
	}

	index(ids = []) {
		if (!ids.length) return this.#courses;

		const courses = ids.map(id => this.show(id))

		if (!courses.filter(c => !!c).length) return undefined;

		return courses;
	}

}

const port = process.env.PORT || 3000;

function createCurso(id) {
	const titulo = faker.name.jobArea();
	return {
		id: id.toString(),
		titulo,
		aulas: faker.datatype.number(10, 40),
		hot_site: `${ faker.internet.url() }/${ id }/${ titulo }`,
		alunos: `http://localhost:${ port }/course/${ id }/students`,
		_self: `http://localhost:${ port }/course/${ id }`,
	}
}

module.exports = {
	courses: new Courses(),
	createCurso
}