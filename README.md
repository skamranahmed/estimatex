## 👨‍💻 EstimateX (`estimatex`)

`estimatex` is a Command-Line Interface (CLI) tool designed to facilitate **story point estimation** of tasks/tickets directly from the terminal.

This client communicates with the [`estimatex-server`](https://github.com/skamranahmed/estimatex-server) via a **WebSocket connection**.

#### 🚀 Demo

##### Room Create Demo
<img alt="Room Create Demo" src="./room-create-demo.gif" />

##### Room Join Demo
<img alt="Room join Demo" src="./room-join-demo.gif" />

### 🙌 Why I Built This
In my team, story point estimation is an integral part of our sprint planning process. It's a collaborative effort that helps us gauge the complexity of tasks and plan our work effectively.

While we often use open-source web based tools for this, I wondered if I could create a CLI-based tool for the same process.

That curiosity led me to build `estimatex`, a cli-based tool for story point estimation.

The motivation was to:

- Explore WebSocket-based real-time communication.

- Create a lightweight, terminal-first alternative for story point estimation sessions.

### ✨ Features

#### Room Management
- **Create Estimation Rooms**: Create private rooms for estimation sessions with configurable maximum capacity
- **Join Existing Rooms**: Join ongoing estimation sessions using room IDs
- **Real-time Updates**: Get real-time updates in your terminal session when team members join the room

#### Estimation Process
- **Story Point Voting**: Choose from standard Fibonacci sequence values (1, 2, 3, 5, 8, 13, 21) for estimations
- **Anonymous Voting**: Votes remain hidden until all team members have submitted their estimates
- **Vote Revelation**: Room admin can trigger vote revelation once all members have voted
- **Results Display**: Clear tabulated display of voting results showing:
  - Distribution of votes
  - Individual member votes
  - Vote count per story point value

#### User Experience
- **Interactive CLI**: Simple command-line interface with clear prompts and emoji-enhanced feedback
- **Secure WebSocket Communication**: Uses WSS (WebSocket Secure) protocol in production for secure real-time communication
- **Graceful Exit**: Clean connection termination with CTRL+C, ensuring proper cleanup of WebSocket connections

#### Administrative Controls
- **Room Capacity Management**: Automatic handling of room capacity limits
- **Session Control**: Room admin has control over when to start voting and reveal results
- **Development Mode**: Built-in development mode for testing with local server setup

### 🛠️ Installation & Usage

#### Option 1: Download Binary (Recommended)
Pre-built binaries are available for Linux, macOS, and Windows. Download the latest release from the [releases page](https://github.com/skamranahmed/estimatex/releases).

1. **For macOS/Linux:**
   ```bash
   # Download the appropriate tar.gz for your OS
   tar xvf estimatex_<OS>_<ARCH>.tar.gz

   # Make the binary executable
   chmod +x ./estimatex

   # Move binary to a directory in your PATH
   sudo mv estimatex /usr/local/bin/
   ```

2. **For Windows:**
   - Download the `.zip` file for Windows
   - Extract the archive
   - Add the extracted directory to your PATH or run `estimatex.exe` directly

#### Option 2: Run from Source
If you prefer to run from source, you'll need Go installed on your system.

1. Clone the repository:
   ```bash
   git clone https://github.com/skamranahmed/estimatex.git
   cd estimatex
   ```

2. Install dependencies:
   ```bash
   make dep
   ```

3. Run the application:
   ```bash
   make run
   ```

#### Development Mode
By default, the application runs in production mode connecting to the hosted server. To run in development mode:

1. Set `isDevelopment = true` in `main.go`
2. Ensure you have the [estimatex-server](https://github.com/skamranahmed/estimatex-server) running locally
3. Run using `make run`

### 📝 License
This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/)