## EstimateX (`estimatex`)

`estimatex` is a Command-Line Interface (CLI) tool designed to facilitate **story point estimation** of tasks/tickets directly from the terminal.

This client communicates with the `estimatex-server` via a **WebSocket connection**.

#### Demo

##### Room Create Demo
<img alt="Room Create Demo" src="./room-create-demo.gif" />

##### Room Join Demo
<img alt="Room join Demo" src="./room-join-demo.gif" />

### Why I Built This
In my team, story point estimation is an integral part of our sprint planning process. It's a collaborative effort that helps us gauge the complexity of tasks and plan our work effectively.

While we often use open-source web based tools for this, I wondered if I could create a CLI-based tool for the same process.

That curiosity led me to build `estimatex`, a cli-based tool for story point estimation.

The motivation was to:

- Explore WebSocket-based real-time communication.

- Create a lightweight, terminal-first alternative for story point estimation sessions.

### Todo 📝

- [ ] Integrate goreleaser
- [ ] Explain the features

### License
This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/)