const API_URL = "http://localhost:8080";
let currentUser = null;
let currentTodoId = null; // For comments

document.addEventListener("DOMContentLoaded", init);

async function init() {
    // 1. Fetch Users for the dropdown
    const res = await fetch(`${API_URL}/users`);
    const users = await res.json();
    
    const select = document.getElementById("userSelect");
    const assignSelect = document.getElementById("assignSelect");

    users.forEach(u => {
        // Login Dropdown
        const opt = document.createElement("option");
        opt.value = u.id;
        opt.text = `${u.username} (${u.role})`;
        opt.dataset.role = u.role;
        select.appendChild(opt);

        // Assignment Dropdown (For admin)
        const assignOpt = document.createElement("option");
        assignOpt.value = u.id;
        assignOpt.text = u.username;
        assignSelect.appendChild(assignOpt);
    });

    loadApp();
}

function loadApp() {
    const select = document.getElementById("userSelect");
    const selectedOpt = select.options[select.selectedIndex];
    
    currentUser = {
        id: parseInt(select.value),
        role: selectedOpt.dataset.role
    };

    // Show/Hide Admin Panel
    document.getElementById("adminPanel").style.display = 
        (currentUser.role === 'admin') ? 'flex' : 'none';

    fetchTodos();
}

async function fetchTodos() {
    // Pass User ID in header
    const res = await fetch(`${API_URL}/todos`, {
        headers: { "X-User-ID": currentUser.id }
    });
    const todos = await res.json() || [];
    
    const list = document.getElementById("todoList");
    list.innerHTML = "";

    todos.forEach(t => {
        const li = document.createElement("li");
        li.innerHTML = `
            <div>
                <strong>${t.title}</strong> <br>
                <small>Assigned to: ${t.assigned_name || 'Unassigned'}</small>
            </div>
            <div>
                <button onclick="openComments(${t.id})">💬 Comments</button>
                <button onclick="toggleTodo(${t.id}, ${!t.completed})">
                    ${t.completed ? "✅ Done" : "⬜ To Do"}
                </button>
            </div>
        `;
        list.appendChild(li);
    });
}

async function addTodo() {
    const title = document.getElementById("todoInput").value;
    const assignedTo = document.getElementById("assignSelect").value;

    if(!title || !assignedTo) return alert("Fill all fields");

    await fetch(`${API_URL}/todos`, {
        method: "POST",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ title, assigned_to: parseInt(assignedTo) })
    });
    fetchTodos();
}

async function toggleTodo(id, status) {
    await fetch(`${API_URL}/todos/${id}`, {
        method: "PUT",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ completed: status })
    });
    fetchTodos();
}

// --- Comments Logic ---

async function openComments(id) {
    currentTodoId = id;
    const dialog = document.getElementById("commentDialog");
    const list = document.getElementById("commentsList");
    list.innerHTML = "Loading...";
    
    dialog.showModal();

    const res = await fetch(`${API_URL}/todos/${id}/comments`, {
        headers: { "X-User-ID": currentUser.id }
    });
    const comments = await res.json() || [];

    list.innerHTML = comments.map(c => `
        <p><strong>${c.username}:</strong> ${c.content}</p>
    `).join("");
}

async function postComment() {
    const input = document.getElementById("commentInput");
    if(!input.value) return;

    await fetch(`${API_URL}/todos/${currentTodoId}/comments`, {
        method: "POST",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ content: input.value })
    });
    
    input.value = "";
    openComments(currentTodoId); // Refresh list
}