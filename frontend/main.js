// Listen for text updates from Go backend
window.runtime.EventsOn("updateText", (data) => {
    const status = document.getElementById("status");
    status.textContent = data.text;
    status.className = data.state || "waiting";
});

// Signal that frontend is ready
window.runtime.EventsEmit("frontendReady");
