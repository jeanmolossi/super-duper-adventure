module.exports = {
	sleep: async (sleepTimeInMs = 25) => {
		await new Promise(resolve => setTimeout(resolve, sleepTimeInMs)).catch(console.warn)
	}
}