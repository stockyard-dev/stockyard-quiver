package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Quiver</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono);cursor:pointer;transition:border-color .2s}
.st:hover,.st.active{border-color:var(--rust)}.st.active .st-v{color:var(--rust)}
.st-v{font-size:1.3rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:200px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-p:hover{background:#d4682f}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.item:hover{border-color:var(--leather)}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.item-title{font-size:.9rem;margin-bottom:.15rem}.item-title a{color:var(--cream);text-decoration:none}.item-title a:hover{color:var(--rust)}
.item-author{font-size:.75rem;color:var(--cd);font-style:italic}
.item-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.type-badge{font-family:var(--mono);font-size:.5rem;padding:.15rem .4rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}
.type-badge.article{border-color:var(--blue);color:var(--blue)}
.type-badge.book{border-color:var(--gold);color:var(--gold)}
.type-badge.paper{border-color:var(--green);color:var(--green)}
.type-badge.video{border-color:var(--red);color:var(--red)}
.type-badge.course{border-color:var(--leather);color:var(--leather)}
.type-badge.podcast{border-color:var(--rust);color:var(--rust)}
.stars{color:var(--gold);letter-spacing:2px;font-size:.7rem}
.tag{font-family:var(--mono);font-size:.5rem;padding:.1rem .35rem;background:var(--bg3);color:var(--cm)}
.item-notes{font-size:.75rem;color:var(--cm);margin-top:.4rem;padding:.4rem .6rem;border-left:2px solid var(--bg3);font-style:italic}
.status-pill{font-family:var(--mono);font-size:.5rem;padding:.15rem .4rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.status-pill.unread{border-color:var(--cm);color:var(--cm)}
.status-pill.reading{border-color:var(--gold);color:var(--gold)}
.status-pill.finished{border-color:var(--green);color:var(--green)}
.item-actions{display:flex;gap:.3rem;flex-shrink:0}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr textarea{resize:vertical}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.rating-input{display:flex;gap:.2rem;margin-top:.2rem}.rating-input span{cursor:pointer;font-size:1.1rem;color:var(--bg3);transition:color .15s}
.rating-input span.on{color:var(--gold)}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
.count-label{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-bottom:.5rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.toolbar{flex-direction:column}.search{min-width:100%}.row2{grid-template-columns:1fr}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> QUIVER</h1><button class="btn btn-p" onclick="openForm()">+ Add Item</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" type="text" placeholder="Search titles, authors, tags..." oninput="render()">
<select class="filter-sel" id="type-filter" onchange="render()"><option value="">All Types</option><option value="article">Article</option><option value="book">Book</option><option value="paper">Paper</option><option value="course">Course</option><option value="video">Video</option><option value="podcast">Podcast</option></select>
</div>
<div class="count-label" id="count"></div>
<div id="items"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],filter='all',editId=null;
async function load(){var r=await fetch(A+'/reading_items').then(function(r){return r.json()});items=r.items||r.reading_items||[];renderStats();render();}
function renderStats(){var total=items.length,unread=items.filter(function(i){return i.status==='unread'}).length,reading=items.filter(function(i){return i.status==='reading'}).length,done=items.filter(function(i){return i.status==='finished'}).length;
document.getElementById('stats').innerHTML=[{l:'Total',v:total,f:'all'},{l:'To Read',v:unread,f:'unread'},{l:'Reading',v:reading,f:'reading'},{l:'Finished',v:done,f:'finished'}].map(function(x){return '<div class="st'+(filter===x.f?' active':'')+'" onclick="setFilter(\''+x.f+'\')"><div class="st-v">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');}
function setFilter(f){filter=f;renderStats();render();}
function render(){var q=(document.getElementById('search').value||'').toLowerCase(),tf=document.getElementById('type-filter').value,f=items;
if(filter!=='all')f=f.filter(function(i){return i.status===filter});
if(tf)f=f.filter(function(i){return i.type===tf});
if(q)f=f.filter(function(i){return(i.title||'').toLowerCase().includes(q)||(i.author||'').toLowerCase().includes(q)||(i.tags||'').toLowerCase().includes(q)});
document.getElementById('count').textContent=f.length+' item'+(f.length!==1?'s':'');
if(!f.length){document.getElementById('items').innerHTML='<div class="empty">No items found.</div>';return;}
var h='';f.forEach(function(i){var hasUrl=i.url&&i.url.trim();
h+='<div class="item"><div class="item-top"><div>';
h+='<div class="item-title">'+(hasUrl?'<a href="'+esc(i.url)+'" target="_blank" rel="noopener">'+esc(i.title)+' &#8599;</a>':esc(i.title))+'</div>';
if(i.author)h+='<div class="item-author">by '+esc(i.author)+'</div>';
h+='</div><div class="item-actions">';
if(i.status==='unread')h+='<button class="btn btn-sm" onclick="setStatus(\''+i.id+'\',\'reading\')">Start</button>';
if(i.status==='reading')h+='<button class="btn btn-sm" onclick="setStatus(\''+i.id+'\',\'finished\')">Done</button>';
if(i.status==='finished')h+='<button class="btn btn-sm" onclick="setStatus(\''+i.id+'\',\'unread\')">Re-read</button>';
h+='<button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button>';
h+='</div></div><div class="item-meta">';
h+='<span class="type-badge '+i.type+'">'+i.type+'</span>';
h+='<span class="status-pill '+i.status+'">'+i.status+'</span>';
if(i.rating&&i.rating>0){var s='';for(var x=0;x<5;x++)s+=(x<i.rating?'&#9733;':'&#9734;');h+='<span class="stars">'+s+'</span>';}
if(i.tags){i.tags.split(',').forEach(function(t){t=t.trim();if(t)h+='<span class="tag">#'+esc(t)+'</span>';});}
h+='</div>';
if(i.notes)h+='<div class="item-notes">'+esc(i.notes)+'</div>';
h+='</div>';});
document.getElementById('items').innerHTML=h;}
async function setStatus(id,status){var body={status:status};if(status==='finished')body.completed_at=new Date().toISOString();if(status==='unread')body.completed_at='';await fetch(A+'/reading_items/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});load();}
async function del(id){if(!confirm('Remove this item?'))return;await fetch(A+'/reading_items/'+id,{method:'DELETE'});load();}
function formHTML(item){var i=item||{title:'',author:'',url:'',type:'article',status:'unread',rating:0,tags:'',notes:''};var isEdit=!!item;
var types=['article','book','paper','course','video','podcast'];var statuses=['unread','reading','finished'];
var h='<h2>'+(isEdit?'EDIT ITEM':'ADD TO READING LIST')+'</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'" placeholder="e.g. Designing Data-Intensive Applications"></div>';
h+='<div class="row2"><div class="fr"><label>Author</label><input id="f-author" value="'+esc(i.author)+'"></div>';
h+='<div class="fr"><label>URL</label><input id="f-url" value="'+esc(i.url)+'" placeholder="https://"></div></div>';
h+='<div class="row2"><div class="fr"><label>Type</label><select id="f-type">';
types.forEach(function(t){h+='<option value="'+t+'"'+(i.type===t?' selected':'')+'>'+t.charAt(0).toUpperCase()+t.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Status</label><select id="f-status">';
statuses.forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div></div>';
h+='<div class="fr"><label>Rating</label><div class="rating-input" id="f-rating" data-val="'+(i.rating||0)+'">';
for(var n=1;n<=5;n++){h+='<span class="'+(n<=(i.rating||0)?'on':'')+'" onclick="setRating('+n+')">&#9733;</span>';}
h+='</div></div>';
h+='<div class="fr"><label>Tags</label><input id="f-tags" value="'+esc(i.tags)+'" placeholder="comma separated"></div>';
h+='<div class="fr"><label>Notes</label><textarea id="f-notes" rows="3" placeholder="Your thoughts...">'+esc(i.notes)+'</textarea></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';
return h;}
function setRating(n){var el=document.getElementById('f-rating');var cur=parseInt(el.dataset.val)||0;var val=(cur===n)?0:n;el.dataset.val=val;var spans=el.querySelectorAll('span');for(var j=0;j<spans.length;j++){spans[j].className=(j<val)?'on':'';}}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-title').focus();}
function openEdit(id){var item=null;for(var j=0;j<items.length;j++){if(items[j].id===id){item=items[j];break;}}if(!item)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(item);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var title=document.getElementById('f-title').value.trim();if(!title){alert('Title is required');return;}
var body={title:title,author:document.getElementById('f-author').value.trim(),url:document.getElementById('f-url').value.trim(),type:document.getElementById('f-type').value,status:document.getElementById('f-status').value,rating:parseInt(document.getElementById('f-rating').dataset.val)||0,tags:document.getElementById('f-tags').value.trim(),notes:document.getElementById('f-notes').value.trim()};
if(body.status==='finished'&&!editId)body.completed_at=new Date().toISOString();
if(editId){await fetch(A+'/reading_items/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/reading_items',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
