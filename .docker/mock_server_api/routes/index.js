const { Router } = require("express");
const studentsResource = require("./students");
const coursesResource = require("./courses");
const seedResource = require("./seed");

module.exports = function makeRoutes() {
	const router = Router();

	router.use(studentsResource);
	router.use(coursesResource);
	router.use('/seed', seedResource)

	return router;
}