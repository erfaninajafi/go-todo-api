const API_URL = "http://localhost:8080/todos";

// Load todos on page start
document.addEventListener("DOMContentLoaded", fetchTodos);

async function fetchTodos() {
    const response = await fetch(API_URL);
    const todos = await response.json();
    const list = document.getElementById("todoList");
    list.innerHTML = "";

    // Handle case where DB is empty (returns null/undefined sometimes)
    if (!todos) return;

    todos.forEach(todo => {
        const li = document.createElement("li");
        if (todo.completed) li.classList.add("completed");

        li.innerHTML = `
            <span onclick="toggleTodo(${todo.id}, ${!todo.completed})" style="cursor:pointer;">
                ${todo.completed ? "✅" : "⬜"} ${todo.title}
            </span>
            <button class="delete-btn" onclick="deleteTodo(${todo.id})">Delete</button>
        `;
        list.appendChild(li);
    });
}

async function addTodo() {
    const input = document.getElementById("todoInput");
    const title = input.value.trim();
    if (!title) return;

    await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: title })
    });

    input.value = "";
    fetchTodos();
}

async function toggleTodo(id, status) {
    await fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ completed: status })
    });
    fetchTodos();
}

async function deleteTodo(id) {
    if(!confirm("Are you sure?")) return;
    
    await fetch(`${API_URL}/${id}`, {
        method: "DELETE"
    });
    fetchTodos();
}