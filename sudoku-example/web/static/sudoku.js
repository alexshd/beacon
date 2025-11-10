// Sudoku Board HTMX Integration
// Clean JavaScript - no inline handlers!

// Current selected cell for number placement
let selectedCell = null;

// Initialize board interactions
document.addEventListener('DOMContentLoaded', () =>
{
  setupCellSelection();
  setupKeyboardInput();
  refreshBoard();
});

// Setup cell click handlers
function setupCellSelection()
{
  document.addEventListener('click', (e) =>
  {
    if (e.target.classList.contains('sudoku-cell'))
    {
      // Remove previous selection
      document.querySelectorAll('.sudoku-cell').forEach(cell =>
      {
        cell.classList.remove('selected');
      });

      // Select new cell
      selectedCell = e.target;
      selectedCell.classList.add('selected');
    }
  });
}

// Setup keyboard number input
function setupKeyboardInput()
{
  document.addEventListener('keydown', (e) =>
  {
    if (!selectedCell) return;

    const num = parseInt(e.key);
    if (num >= 1 && num <= 9)
    {
      placeNumber(selectedCell, num);
    } else if (e.key === 'Backspace' || e.key === 'Delete' || e.key === '0')
    {
      placeNumber(selectedCell, 0);
    }
  });
}

// Place number via HTMX
function placeNumber(cell, num)
{
  const row = parseInt(cell.dataset.row);
  const col = parseInt(cell.dataset.col);

  fetch('/place', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ row, col, num })
  })
    .then(response => response.json())
    .then(data =>
    {
      if (data.success)
      {
        refreshBoard();
      }
    })
    .catch(err => console.error('Failed to place number:', err));
}

// Refresh board display via HTMX
function refreshBoard()
{
  htmx.ajax('GET', '/board-html', {
    target: '#sudoku-board',
    swap: 'innerHTML'
  });

  htmx.ajax('GET', '/stats-html', {
    target: '#stats',
    swap: 'innerHTML'
  });
}

// Merge from another server
function mergeFrom(url)
{
  const statusEl = document.getElementById('merge-status');
  statusEl.textContent = 'Merging...';
  statusEl.className = '';

  fetch(url + '/export')
    .then(response => response.json())
    .then(remoteState =>
    {
      return fetch('/merge', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(remoteState)
      });
    })
    .then(response => response.json())
    .then(data =>
    {
      if (data.success)
      {
        statusEl.textContent = `✅ ${data.message} - Merged ${data.filled} cells`;
        statusEl.className = 'success';
        refreshBoard();
      } else
      {
        throw new Error('Merge failed');
      }
    })
    .catch(err =>
    {
      statusEl.textContent = `❌ Merge failed: ${err.message}`;
      statusEl.className = 'error';
    });
}

// Auto-refresh stats every 2 seconds
setInterval(() =>
{
  htmx.ajax('GET', '/stats-html', {
    target: '#stats',
    swap: 'innerHTML'
  });
}, 2000);

// HTMX event handlers
document.body.addEventListener('htmx:afterSwap', (event) =>
{
  if (event.detail.target.id === 'sudoku-board')
  {
    setupCellSelection();
  }
});
