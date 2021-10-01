const amqplib = require('amqplib');
const express = require('express');
const faker = require('faker');
const _ = require('lodash');

const port = process.env.PORT || 3000;

function createUser(curso) {
    return {
        id: faker.datatype.uuid(),
        nome: faker.name.findName(),
        plano: faker.random.arrayElement([400, 300, 200, 100]),
        email: faker.internet.email().toLowerCase(),
        curso,
    }
}

function createCurso(id) {
    const titulo = faker.name.jobArea();
    return {
        id,
        titulo,
        aulas: faker.datatype.number(10, 40),
        alunos: `http://localhost:${port}/curso/${id}/alunos`,
        link: `${faker.internet.url()}/${id}/${titulo}`,
        _self: `http://localhost:${port}/curso/${id}/ver`,
    }
}

let openRabbit;
const channel = rabbitConnect().then(conn => conn.createChannel())

function rabbitConnect() {
    if (!openRabbit) {
        openRabbit = amqplib.connect('amqp://rabbitmq:rabbitmq@gsr_rabbit?channelMax=65535')
    }

    return openRabbit
}

async function publishMessage(id) {
    const queueName = 'queueCourseId';

    return channel
        .then(function (channel) {
            return channel.assertQueue(queueName).then(function () {
                return channel.sendToQueue(queueName, Buffer.from(id), {routingKey: 'course'})
            })
                .catch(console.warn)
        })
}

const app = express();
app.use(express.json());

app.get('/', (req, res) => {
    return res.json({healthy: true})
})

app.get('/cursos', (req, res) => {
    const cursos = [];

    const totalCursos = Math.round(Math.random() * 30) + 1;

    for (let i = 0; i < totalCursos; i++) {
        cursos.push(createCurso(faker.datatype.number(i + 100000, totalCursos + 100001)))
    }

    const data = {
        totalCursos,
        cursos,
    }

    return res.json(data);
})

app.get('/curso/:id/ver', (req, res) => {
    const {id} = req.params;
    const h = req.headers;

    if (!h['origin-app']) {
        publishMessage(id)
    }

    const failChance = Math.random() * 100

    if (failChance <= 1) {
        return res.status(504).json()
    }

    return res.json(createCurso(id))
})

app.get('/curso/:id/alunos', (req, res) => {

    const {id} = req.params;

    if (!id || isNaN(Number(id))) {
        return res.status(403).json({
            message: `Id do curso nao corresponde a um valor valido`
        })
    }

    const totalUsers = Math.round(Math.random() * 500) + 1;
    const users = [];

    for (let i = 0; i < totalUsers; i++) {
        users.push(createUser(id))
    }

    const data = {
        curso: +req.params.id,
        total: totalUsers,
        users
    }

    return res.json(data);
});

app.get('/fill-queue', async (req, res) => {
    const {query} = req

    let records = 1800
    if (query.records) {
        records = query.records < 65535 ? query.records : 65534
    }

    const idsToFill = []

    for (let i = 0; i < records; i++) {
        const id = faker.datatype.number(i + 100000, records + 100001).toString()
        idsToFill.push(id)
    }

    const promises = idsToFill.map(async id => {
        console.log("Generating records...", id)
        return publishMessage(id)
    })

    await Promise.all(_.chunk(promises, 250))

    console.log(`Generated ${promises.length} records`)

    return res.json({
        message: `Queue filled with ${records}`
    })
})

app.listen(port, () => console.log(`Running on port ${port}`))