function showToast(message, ms = 2400) {
  const t = document.getElementById('toast');
  t.textContent = message;
  t.classList.add('show');
  clearTimeout(t._timeout);
  t._timeout = setTimeout(() => t.classList.remove('show'), ms);
}

async function requestJSON(path, opts = {}){
  const res = await fetch(path, Object.assign({headers:{'Content-Type':'application/json'}}, opts));
  if (!res.ok) {
    const text = await res.text().catch(()=>res.statusText||'error');
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.status !== 204 ? res.json() : null;
}

function makeCell(text){ const td = document.createElement('td'); td.textContent = text; return td }

function emptyRow() {
  const tbody = document.querySelector('#entries tbody');
  tbody.innerHTML = `<tr class="empty"><td colspan="4">No entries yet — add your first weight.</td></tr>`;
}

async function listWeights(){
  let data;
  try{ data = await requestJSON('/weights') }catch(e){ console.error(e); showToast('Could not load entries'); return }
  const tbody = document.querySelector('#entries tbody');
  tbody.innerHTML = '';
  if (!data || data.length === 0) return emptyRow();
  data.forEach(entry => {
    const tr = document.createElement('tr');
    tr.appendChild(makeCell(entry.date));
    tr.appendChild(makeCell(String(entry.weight)));
    tr.appendChild(makeCell(entry.notes || ''));
    const actions = document.createElement('td');
    const edit = document.createElement('button');
    edit.className = 'action-btn edit';
    edit.textContent = 'Edit';
    edit.onclick = () => startEdit(entry);
    const del = document.createElement('button');
    del.className = 'action-btn delete';
    del.textContent = 'Delete';
    del.onclick = () => confirmDelete(entry.id);
    actions.appendChild(edit);
    actions.appendChild(document.createTextNode(' '));
    actions.appendChild(del);
    tr.appendChild(actions);
    tbody.appendChild(tr);
  });
}

async function createWeight(evt){
  evt && evt.preventDefault && evt.preventDefault();
  const date = document.getElementById('date').value;
  const weight = parseFloat(document.getElementById('weight').value);
  const notes = document.getElementById('notes').value;
  if (!date || Number.isNaN(weight)) return showToast('Please provide date and weight');
  try{
    await requestJSON('/weights', { method:'POST', body: JSON.stringify({ date, weight, notes }) });
    document.getElementById('createForm').reset();
    listWeights();
    showToast('Entry added');
  }catch(e){ console.error(e); showToast('Failed to add entry') }
}

async function confirmDelete(id){
  if (!confirm('Delete this entry?')) return;
  try{ await requestJSON('/weights/'+id, { method: 'DELETE' }); showToast('Deleted'); listWeights() }catch(e){ console.error(e); showToast('Delete failed') }
}

function startEdit(entry){
  // populate form with existing values - reuse create form for simplicity
  document.getElementById('date').value = entry.date;
  document.getElementById('weight').value = entry.weight;
  document.getElementById('notes').value = entry.notes || '';
  const form = document.getElementById('createForm');
  form.dataset.editId = entry.id;
  document.querySelector('.form-actions .primary').textContent = 'Save changes';
  showToast('Editing entry — submit to save');
}

async function submitEdit(form){
  const id = form.dataset.editId;
  if (!id) return createWeight();
  const date = document.getElementById('date').value;
  const weight = parseFloat(document.getElementById('weight').value);
  const notes = document.getElementById('notes').value;
  try{
    await requestJSON('/weights/'+id, { method:'PUT', body: JSON.stringify({ date, weight, notes }) });
    delete form.dataset.editId;
    form.reset();
    document.querySelector('.form-actions .primary').textContent = 'Add entry';
    showToast('Saved');
    listWeights();
  }catch(e){ console.error(e); showToast('Save failed') }
}

document.addEventListener('DOMContentLoaded', ()=>{
  listWeights();
  const form = document.getElementById('createForm');
  form.addEventListener('submit', (e)=>{ e.preventDefault(); submitEdit(form) });
  document.getElementById('refreshBtn').addEventListener('click', listWeights);
});
