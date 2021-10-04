const _ = require("lodash");
const faker = require("faker");

class Students {

	#users = [];

	create(payload) {
		const { id, nome, plano, email, curso } = payload;

		this.#users.push(payload);

		return {
			id,
			nome,
			plano,
			email,
			curso,
		}
	}

	show(id) {
		if (!id) return;

		const user = _.find(this.#users, { id });

		return user || {
			id,
			code: 404,
			message: "User not found"
		};
	}

	index(ids = []) {
		if (!ids.length) {
			return this.#users
		}

		const users = ids.map(id => this.show(id))

		if (!totalResults(users)) {
			return undefined;
		}

		return users;
	}

	indexFromCourse(courseId) {
		if (!courseId) return [];

		const users = _.filter(this.#users, { curso: courseId.toString() })

		if (!totalResults(users)) {
			return [];
		}

		return users;
	}

}

function totalResults(r = []) {
	return r.filter(i => !!i).length
}

const port = process.env.PORT || 3000

function createUser(curso) {
	const id = faker.datatype.uuid();
	return {
		id,
		nome: faker.name.findName(),
		plano: faker.random.arrayElement([400, 300, 200, 100]),
		email: faker.internet.email().toLowerCase(),
		curso: curso.toString(),
		_self: `http://localhost:${ port }/student/${ id }`
	}
}

module.exports = {
	students: new Students(),
	createUser
}