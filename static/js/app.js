let username;
let socket
function submitUsername() {
    username = document.getElementById('usernameInput').value;

    if (username) {
        document.getElementById('usernameForm').style.display = 'none';
        document.getElementById('chatContainer').style.display = 'block';
        document.getElementById('usernameDisplay').innerText = username;

        socket = new WebSocket('ws://localhost:8080/ws');

        socket.addEventListener('open', () => {
            const message = {
                type: 'username',
                content: username
            };
            socket.send(JSON.stringify(message));
        });

        socket.addEventListener('message', (event) => {
            console.log(event.data)
            const message = JSON.parse(event.data);
            console.log(message)
            if (message.type === 'message') {
                const content = JSON.parse(message.content)
                console.log(content)
               const chatDiv = document.getElementById('chat');
                chatDiv.innerHTML += `<p><strong>${content.username}:</strong> ${content.message}</p>`;

            } else if (message.type === 'user_registration') {
                console.log("You registered successfully as " + message.content)
            }
        });

        socket.addEventListener('close', (event) => {
            if (event.code !== 1000) {
                alert('Connection closed unexpectedly. Please refresh the page.');
            }
        });
    }
}

function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const messageContent = messageInput.value;

    if (messageContent) {
        //const socket = new WebSocket('ws://localhost:8080/ws');
        const message = {
            type: 'message',
            content: {
                username,
                message: messageContent
            }
        };
        socket.send(JSON.stringify(message));
        messageInput.value = '';
    }
}
