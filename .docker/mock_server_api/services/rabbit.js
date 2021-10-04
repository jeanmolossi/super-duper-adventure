const amqplib = require("amqplib");

class RabbitMQ {
	#openRabbit;
	#queueName;
	#channel;
	routingKey = 'course';

	constructor(queueName = "queueCourseId") {
		const connection = amqplib.connect('amqp://rabbitmq:rabbitmq@gsr_rabbit?channelMax=65535');

		this.#queueName = queueName
		this.#openRabbit = connection
		this.#channel = connection.then(conn => conn.createChannel())
	}

	#getChannel(payload) {
		return this.#channel.then(function (channel) {
			return { channel, payload };
		});
	}

	#getQueue = ({ channel, payload }) => {
		return channel.assertQueue(this.#queueName)
			.then(function () {
				return { channel, payload }
			})
	}

	#sendToQueue = ({ channel, payload }) => {
		return channel.sendToQueue(this.#queueName, Buffer.from(payload), { routingKey: this.routingKey })
	}

	publish = (payload) => {
		return this.#getChannel(payload)
			.then(this.#getQueue)
			.then(this.#sendToQueue)
	}
}

module.exports = {
	rabbitMQ: new RabbitMQ()
}