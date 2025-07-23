window.onload = function() {
    const canvas = document.getElementById('paint-canvas');
    const pencilBtn = document.getElementById('pencil-btn');
    const eraserBtn = document.getElementById('eraser-btn');
    const colorPicker = document.getElementById('color-picker');

    const ctx = canvas.getContext('2d');
    const ERASER_SIZE = 12;

    let isPainting = false;
    let lastX = 0;
    let lastY = 0;
    let currentColor = colorPicker.value;
    let currentTool = 'pencil';

    const pathParts = window.location.pathname.split('/');
    const roomID = pathParts[pathParts.length - 1];

    const url = "ws://" + window.location.host + "/ws/" + roomID;
    const conn = new WebSocket(url);

    console.log("Подключение к комнате:", roomID);
    console.log("URL для WebSocket:", url);

    const userCountSpan = document.getElementById('user-count');

    conn.onopen = function() {
        console.log("Соединение установлено.");
    };

    conn.onclose = function () {
        console.log("Соединение закрыто.");
    };

    conn.onmessage = function (evt) {
        const data = JSON.parse(evt.data);
        if (data.type === 'user_count') {
            userCountSpan.textContent = data.count;
        } else {
            drawOnCanvas(data);
        }
    };

    function updateCursor() {
        if (currentTool === 'pencil') {
            canvas.style.cursor = 'crosshair';
        } else if (currentTool === 'eraser') {
            const eraserCursor = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgdmlld0JveD0iMCAwIDI1IDI1Ij48cmVjdCB4PSIwIiB5PSIwIiB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGZpbGw9IndoaXRlIiBzdHJva2U9ImJsYWNrIiBzdHJva2Utd2lkdGg9IjEiLz48L3N2Zz4=";
            canvas.style.cursor = `url('${eraserCursor}') 12 12, auto`;

        }
    }

    function drawOnCanvas(data) {
        ctx.beginPath();
        ctx.strokeStyle = data.color
        ctx.lineWidth = data.lineWidth
        ctx.lineCap = data.lineCap
        ctx.moveTo(data.x1, data.y1);
        ctx.lineTo(data.x2, data.y2);
        ctx.stroke();
   }

    function draw(e) {
        if (!isPainting) {
            return;
        }

        const data = {
            type: "draw",
            x1: lastX,
            y1: lastY,
            x2: e.offsetX,
            y2: e.offsetY,
            color: (currentTool === 'eraser') ? '#FFFFFF' : currentColor,
            lineWidth: (currentTool === 'eraser') ? ERASER_SIZE : 5,
            lineCap: (currentTool === 'eraser') ? 'square' : 'round',
        };

        console.log("color to draw", data.color)

        drawOnCanvas(data);

        conn.send(JSON.stringify(data));

        [lastX, lastY] = [e.offsetX, e.offsetY];
    }

    canvas.addEventListener('mousedown', (e) => {
        isPainting = true;
        [lastX, lastY] = [e.offsetX, e.offsetY];
    });

    canvas.addEventListener('mousemove', draw);

    canvas.addEventListener('mouseup', () => {
        isPainting = false;
    });

    canvas.addEventListener('mouseout', () => {
        isPainting = false;
    });

    colorPicker.addEventListener('change', (e) => {
        currentColor = e.target.value;
        currentTool = 'pencil';
        console.log("cursor color changed", currentColor)
        updateCursor();
    });

    pencilBtn.addEventListener('click', () => {
        currentTool = 'pencil';
        updateCursor();
    });

    eraserBtn.addEventListener('click', () => {
        currentTool = 'eraser';
        updateCursor();
    });

    updateCursor();
};