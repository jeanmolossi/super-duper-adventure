const { Router } = require("express");
const { students, createUser } = require("../services/students");

const studentsResource = Router();

studentsResource.post('/students', (req, res) => {
	const { id } = req.body;

	const user = students.create(createUser(id))

	return res.json(user)
})

studentsResource.get('/students', (req, res) => {
	const { query } = req;

	const ids = query.ids?.split(',');

	const users = students.index(ids);

	if (!users) {
		return res.status(404).json({ message: "Users not found" })
	}

	return res.json(users)
})

studentsResource.get('/student/:id', (req, res) => {
	const { id } = req.params;

	const user = students.show(id);

	if (!user) {
		return res.status(404).json({ message: 'User not found' })
	}
	return res.json(user)
})

module.exports = studentsResource;