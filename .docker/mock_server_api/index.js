const amqplib = require('amqplib');
const express = require('express');
const faker = require('faker');

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
        aulas: faker.datatype.number(10,40),
        alunos: `http://localhost:3001/curso/${id}/alunos`,
        link: `${faker.internet.url()}/${id}/${titulo}`,
        _self: `http://localhost:3001/curso/${id}/ver`,
    }
}

const app = express();
app.use(express.json());

const port = 3000;

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
    const { id } = req.params;

    const openRabbit = amqplib.connect('amqp://rabbitmq:rabbitmq@gsr_rabbit')
    const queueName = 'queueCourseId';

    openRabbit.then(function (conn) {
        return conn.createChannel()
    })
        .then(function (channel) {
            return channel.assertQueue(queueName).then(function () {
                return channel.sendToQueue(queueName, Buffer.from(id), { routingKey: 'course' })
            })
        })

    return res.json(createCurso(id))
})

app.get('/curso/:id/alunos', (req, res) => {

    const {id} = req.params;

    if (!id || isNaN(Number(id))) {
        return res.status(403).json({
            message: `Id do curso nao corresponde a um valor valido`
        })
    }

    const totalUsers = Math.round(Math.random() * 3000) + 1;
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

app.listen(port, () => console.log(`Running on port ${port}`))