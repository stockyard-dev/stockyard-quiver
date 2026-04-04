package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Quiver</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono)}
.st-v{font-size:1.2rem}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.1rem}
.tabs{display:flex;gap:.3rem;margin-bottom:1rem;font-family:var(--mono)}
.tab{font-size:.65rem;padding:.3rem .7rem;border:1px solid var(--bg3);background:var(--bg);color:var(--cm);cursor:pointer}.tab:hover{border-color:var(--leather)}.tab.active{border-color:var(--rust);color:var(--rust)}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem}
.item-title{font-size:.95rem;margin-bottom:.2rem}
.item-title a{color:var(--cream)}.item-title a:hover{color:var(--rust)}
.item-author{font-size:.78rem;color:var(--cd);font-style:italic}
.item-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.6rem;flex-wrap:wrap}
.type-badge{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}
.stars{color:var(--gold);letter-spacing:1px}
.tag{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cm)}
.item-notes{font-size:.78rem;color:var(--cm);margin-top:.4rem;padding:.4rem;border-left:2px solid var(--bg3);font-style:italic}
.btn{font-family:var(--mono);font-size:.6rem;padding:.2rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.8rem}
</style></head><body>
<div class="hdr"><h1>QUIVER</h1><button class="btn btn-p" onclick="openForm()">+ Add</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="tabs" id="tabs"></div>
<div id="items"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let items=[],filter='unread';
async function load(){const r=await fetch(A+'/items').then(r=>r.json());items=r.items||r.reading_items||[];
const unread=items.filter(i=>i.status==='unread').length;
const reading=items.filter(i=>i.status==='reading').length;
const done=items.filter(i=>i.status==='finished').length;
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+unread+'</div><div class="st-l">To Read</div></div><div class="st"><div class="st-v">'+reading+'</div><div class="st-l">Reading</div></div><div class="st"><div class="st-v">'+done+'</div><div class="st-l">Finished</div></div>';
document.getElementById('tabs').innerHTML=['all','unread','reading','finished'].map(t=>'<button class="tab'+(filter===t?' active':'')+'" onclick="setFilter(\''+t+'\')">'+t+'</button>').join('');
render();}
function setFilter(f){filter=f;render();}
function render(){let filtered=filter==='all'?items:items.filter(i=>i.status===filter);
if(!filtered.length){document.getElementById('items').innerHTML='<div class="empty">'+(filter==='all'?'Your reading list is empty.':'No '+filter+' items.')+'</div>';return;}
let h='';filtered.forEach(i=>{
h+='<div class="item"><div style="display:flex;justify-content:space-between;align-items:flex-start"><div><div class="item-title">'+(i.url?'<a href="'+esc(i.url)+'" target="_blank">'+esc(i.title)+'</a>':esc(i.title))+'</div>';
if(i.author)h+='<div class="item-author">by '+esc(i.author)+'</div>';
h+='</div><div style="display:flex;gap:.3rem">';
if(i.status==='unread')h+='<button class="btn" onclick="setStatus(\''+i.id+'\',\'reading\')">Start</button>';
if(i.status==='reading')h+='<button class="btn" onclick="setStatus(\''+i.id+'\',\'finished\')">Finish</button>';
if(i.status==='finished')h+='<button class="btn" onclick="setStatus(\''+i.id+'\',\'unread\')">Reread</button>';
h+='<button class="btn" onclick="del(\''+i.id+'\')" style="color:var(--cm)">✕</button></div></div>';
h+='<div class="item-meta"><span class="type-badge">'+i.type+'</span>';
if(i.rating){const s="★".repeat(i.rating)+"☆".repeat(5-i.rating);h+='<span class="stars">'+s+'</span>';}
if(i.tags){i.tags.split(',').forEach(t=>{if(t.trim())h+='<span class="tag">'+esc(t.trim())+'</span>';});}
h+='</div>';
if(i.notes)h+='<div class="item-notes">'+esc(i.notes)+'</div>';
h+='</div>';});
document.getElementById('items').innerHTML=h;}
async function setStatus(id,status){await fetch(A+'/items/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({status})});load();}
async function del(id){if(confirm('Remove?')){await fetch(A+'/items/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>Add to Reading List</h2><div class="fr"><label>Title</label><input id="f-t" placeholder="e.g. Designing Data-Intensive Applications"></div><div class="fr"><label>Author</label><input id="f-a"></div><div class="fr"><label>URL</label><input id="f-u" placeholder="https://"></div><div class="fr"><label>Type</label><select id="f-ty"><option value="article">Article</option><option value="book">Book</option><option value="paper">Paper</option><option value="course">Course</option><option value="video">Video</option><option value="podcast">Podcast</option></select></div><div class="fr"><label>Tags</label><input id="f-tg" placeholder="comma separated"></div><div class="fr"><label>Notes</label><textarea id="f-n" rows="2"></textarea></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Add</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/items',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({title:document.getElementById('f-t').value,author:document.getElementById('f-a').value,url:document.getElementById('f-u').value,type:document.getElementById('f-ty').value,tags:document.getElementById('f-tg').value,notes:document.getElementById('f-n').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
