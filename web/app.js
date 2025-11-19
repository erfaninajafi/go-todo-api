const API_URL = "http://localhost:8080";
let isSignup = false;
let currentUser = null; 
let currentTodoId = null;


function toggleAuthMode() {
    isSignup = !isSignup;
    document.getElementById("authTitle").innerText = isSignup ? "Sign Up" : "Login";
    document.getElementById("authBtn").innerText = isSignup ? "Create Account" : "Login";
    document.getElementById("toggleAuthText").innerHTML = isSignup 
        ? 'Have an account? <a href="#" onclick="toggleAuthMode()">Login</a>' 
        : 'Don\'t have an account? <a href="#" onclick="toggleAuthMode()">Sign up</a>';
}

async function handleAuth() {
    const u = document.getElementById("username").value;
    const p = document.getElementById("password").value;
    const endpoint = isSignup ? "/signup" : "/login";

    const res = await fetch(API_URL + endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: u, password: p })
    });

    if (res.ok) {
        const data = await res.json();
        currentUser = isSignup ? {id: data.id, role: data.role, username: u} : data;
        
        document.getElementById("authSection").style.display = "none";
        document.getElementById("dashboardSection").style.display = "block";
        initDashboard();
    } else {
        alert("Auth failed");
    }
}

function logout() {
    currentUser = null;
    location.reload();
}


async function initDashboard() {
    document.getElementById("welcomeMsg").innerText = `Hello, ${currentUser.username} (${currentUser.role})`;

    if (currentUser.role === 'admin') {
        document.getElementById("adminPanel").style.display = "block";
        loadUsersForDropdown();
    } else {
        document.getElementById("adminPanel").style.display = "none";
    }

    loadTodos();
}

async function loadUsersForDropdown() {
    const res = await fetch(`${API_URL}/users`, {
        headers: { "X-User-ID": currentUser.id }
    });
    const users = await res.json();
    const sel = document.getElementById("assignSelect");
    sel.innerHTML = '<option value="">Assign to...</option>';
    users.forEach(u => {
        const opt = document.createElement("option");
        opt.value = u.id;
        opt.text = u.username;
        sel.appendChild(opt);
    });
}

async function loadTodos() {
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
                <strong>${t.title}</strong> 
                <br><small>Assignee: ${t.assigned_name || 'N/A'}</small>
            </div>
            <div>
                <button onclick="openComments(${t.id})">💬 Comment</button>
                <button onclick="toggleTask(${t.id}, ${!t.completed})">
                     ${t.completed ? "✅" : "⬜"}
                </button>
            </div>
        `;
        list.appendChild(li);
    });
}

async function createTask() {
    const title = document.getElementById("todoInput").value;
    const assignedTo = document.getElementById("assignSelect").value;
    
    await fetch(`${API_URL}/todos`, {
        method: "POST",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ title, assigned_to: parseInt(assignedTo) })
    });
    loadTodos();
}

async function toggleTask(id, status) {
    await fetch(`${API_URL}/todos/${id}`, {
        method: "PUT",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ completed: status })
    });
    loadTodos();
}

async function openComments(id) {
    currentTodoId = id;
    document.getElementById("commentDialog").showModal();
    const list = document.getElementById("commentsList");
    list.innerHTML = "Loading...";
    
    if (currentUser.role === 'admin') {
        const res = await fetch(`${API_URL}/todos/${id}/comments`, {
            headers: { "X-User-ID": currentUser.id }
        });
        if (res.status === 403) {
            list.innerHTML = "<i>Permission denied</i>";
        } else {
            const comments = await res.json() || [];
            list.innerHTML = comments.map(c => `<p><b>${c.username}:</b> ${c.content}</p>`).join("");
        }
        document.getElementById("adminOnlyMsg").style.display = "none";
    } else {
        list.innerHTML = "<i>Comments are hidden for users. You can submit a new one below.</i>";
        document.getElementById("adminOnlyMsg").style.display = "block";
    }
}

async function postComment() {
    const content = document.getElementById("commentInput").value;
    await fetch(`${API_URL}/todos/${currentTodoId}/comments`, {
        method: "POST",
        headers: { "X-User-ID": currentUser.id, "Content-Type": "application/json" },
        body: JSON.stringify({ content })
    });
    document.getElementById("commentInput").value = "";
    
    if (currentUser.role === 'admin') {
        openComments(currentTodoId);
    } else {
        alert("Comment sent to admin!");
        document.getElementById("commentDialog").close();
    }
}