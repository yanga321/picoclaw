package channels

// webUIHTML is the full dashboard frontend served at /.
// Backticks inside the template are replaced by the JS template literal syntax
// using String.fromCharCode(96) where needed.
var webUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
<meta name="theme-color" content="#0a0a0f">
<title>Assistant</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0a0a0f;--surface:#111118;--surface2:#1a1a26;--surface3:#222233;
  --border:#2a2a3a;--border2:#3a3a50;
  --text:#e4e4ef;--text2:#8888a0;--text3:#555570;
  --accent:#6c5ce7;--accent2:#a29bfe;--accent-glow:rgba(108,92,231,.15);
  --danger:#ff6b6b;--danger-bg:rgba(255,107,107,.1);
  --success:#2ed573;
  --bot-bg:#1a1a28;
  --radius:14px;--radius-sm:10px;
  --safe-bottom:env(safe-area-inset-bottom,0px);
  --sidebar-w:300px;
  --header-h:60px;
}
html,body{height:100%;width:100%;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;background:var(--bg);color:var(--text);overflow:hidden;-webkit-font-smoothing:antialiased}
.app{display:flex;height:100%;position:relative}

/* Sidebar */
.sidebar{width:var(--sidebar-w);height:100%;background:var(--surface);border-right:1px solid var(--border);display:flex;flex-direction:column;flex-shrink:0;transition:transform .3s cubic-bezier(.4,0,.2,1);z-index:50}
.sidebar-header{padding:16px;display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--border);flex-shrink:0}
.sidebar-logo{width:36px;height:36px;border-radius:10px;background:linear-gradient(135deg,var(--accent),var(--accent2));display:flex;align-items:center;justify-content:center;font-weight:700;font-size:16px;color:#fff;flex-shrink:0;box-shadow:0 4px 16px rgba(108,92,231,.3)}
.sidebar-brand{font-size:16px;font-weight:700;letter-spacing:-.3px;flex:1}
.new-chat-btn{display:flex;align-items:center;gap:8px;margin:12px 12px 8px;padding:10px 14px;background:var(--accent-glow);border:1px solid rgba(108,92,231,.3);border-radius:var(--radius-sm);color:var(--accent2);font-size:13px;font-weight:500;cursor:pointer;transition:all .2s ease}
.new-chat-btn:hover{background:rgba(108,92,231,.2);border-color:var(--accent)}
.new-chat-btn svg{width:16px;height:16px;flex-shrink:0}
.conv-list{flex:1;overflow-y:auto;padding:4px 8px 12px;display:flex;flex-direction:column;gap:2px}
.conv-list::-webkit-scrollbar{width:4px}
.conv-list::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}
.conv-item{display:flex;align-items:center;gap:8px;padding:10px 12px;border-radius:var(--radius-sm);cursor:pointer;transition:all .15s ease;position:relative}
.conv-item:hover{background:var(--surface2)}
.conv-item.active{background:var(--accent-glow);border:1px solid rgba(108,92,231,.2)}
.conv-item:not(.active){border:1px solid transparent}
.conv-title{flex:1;font-size:13px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.conv-time{font-size:11px;color:var(--text3);flex-shrink:0}
.conv-del{display:none;width:24px;height:24px;border-radius:6px;border:none;background:var(--danger-bg);color:var(--danger);cursor:pointer;align-items:center;justify-content:center;font-size:14px;flex-shrink:0}
.conv-item:hover .conv-del{display:flex}
.conv-item:hover .conv-time{display:none}
.sidebar-footer{padding:12px 16px;border-top:1px solid var(--border);display:flex;align-items:center;gap:8px;font-size:12px;color:var(--text3);flex-shrink:0}
.db-badge{padding:2px 8px;border-radius:10px;font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:.5px}
.db-badge.on{background:rgba(46,213,115,.15);color:var(--success)}
.db-badge.off{background:var(--surface2);color:var(--text3)}

/* Main */
.main{flex:1;display:flex;flex-direction:column;min-width:0}
.header{display:flex;align-items:center;gap:12px;padding:0 20px;height:var(--header-h);background:var(--surface);border-bottom:1px solid var(--border);flex-shrink:0}
.menu-btn{display:none;width:36px;height:36px;border-radius:var(--radius-sm);border:1px solid var(--border);background:transparent;color:var(--text);cursor:pointer;align-items:center;justify-content:center;font-size:18px;transition:all .2s;flex-shrink:0}
.menu-btn:hover{background:var(--surface2)}
.header-title{flex:1;font-size:15px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.status-pill{display:flex;align-items:center;gap:6px;padding:4px 12px 4px 8px;border-radius:20px;background:rgba(46,213,115,.08);border:1px solid rgba(46,213,115,.2);font-size:11px;color:var(--success);font-weight:500;flex-shrink:0}
.status-dot{width:6px;height:6px;border-radius:50%;background:var(--success);animation:pulse 2s ease infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.3}}

/* Messages */
.messages{flex:1;overflow-y:auto;overflow-x:hidden;padding:20px 16px;display:flex;flex-direction:column;gap:12px;scroll-behavior:smooth;-webkit-overflow-scrolling:touch}
.messages::-webkit-scrollbar{width:4px}
.messages::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}
.msg{display:flex;gap:10px;max-width:88%;animation:msgIn .35s cubic-bezier(.4,0,.2,1)}
@keyframes msgIn{from{opacity:0;transform:translateY(16px)}to{opacity:1;transform:translateY(0)}}
.msg--user{align-self:flex-end;flex-direction:row-reverse}
.msg--bot{align-self:flex-start}
.msg-avatar{width:30px;height:30px;border-radius:var(--radius-sm);flex-shrink:0;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:600;margin-top:2px}
.msg--bot .msg-avatar{background:var(--surface2);color:var(--accent2)}
.msg--user .msg-avatar{background:linear-gradient(135deg,var(--accent),var(--accent2));color:#fff}
.msg-body{display:flex;flex-direction:column;gap:4px;min-width:0}
.msg-bubble{padding:12px 16px;border-radius:var(--radius);font-size:14px;line-height:1.6;word-wrap:break-word}
.msg--user .msg-bubble{background:linear-gradient(135deg,var(--accent),var(--accent2));color:#fff;border-bottom-right-radius:4px}
.msg--bot .msg-bubble{background:var(--bot-bg);border:1px solid var(--border);border-bottom-left-radius:4px}
.msg-time{font-size:10px;color:var(--text3);padding:0 4px}
.msg--user .msg-time{text-align:right}
.msg--bot .msg-bubble p{margin:0 0 8px}
.msg--bot .msg-bubble p:last-child{margin-bottom:0}
.msg--bot .msg-bubble code{background:var(--surface2);padding:2px 6px;border-radius:4px;font-size:13px;font-family:'SF Mono',Monaco,Consolas,monospace}
.msg--bot .msg-bubble pre{background:var(--surface);border:1px solid var(--border);border-radius:var(--radius-sm);padding:12px 16px;margin:8px 0;overflow-x:auto;position:relative;font-size:13px;line-height:1.5}
.msg--bot .msg-bubble pre code{background:none;padding:0;font-size:13px}
.copy-btn{position:absolute;top:8px;right:8px;padding:4px 10px;border-radius:6px;border:1px solid var(--border);background:var(--surface2);color:var(--text2);font-size:11px;cursor:pointer;opacity:0;transition:opacity .2s}
pre:hover .copy-btn{opacity:1}
.copy-btn:hover{color:var(--accent2);border-color:var(--accent)}

/* Typing */
.typing{display:none;align-self:flex-start;padding:0 16px}
.typing.active{display:flex;align-items:center;gap:8px}
.typing-dots{display:flex;gap:4px}
.typing-dots span{width:6px;height:6px;border-radius:50%;background:var(--text2);animation:bounce 1.4s ease infinite}
.typing-dots span:nth-child(2){animation-delay:.2s}
.typing-dots span:nth-child(3){animation-delay:.4s}
@keyframes bounce{0%,60%,100%{transform:translateY(0)}30%{transform:translateY(-6px)}}
.typing-label{font-size:12px;color:var(--text2)}

/* Welcome */
.welcome{flex:1;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:20px;padding:40px 24px;text-align:center}
.welcome-icon{width:80px;height:80px;border-radius:24px;background:linear-gradient(135deg,var(--accent),var(--accent2));display:flex;align-items:center;justify-content:center;font-size:40px;color:#fff;box-shadow:0 12px 40px rgba(108,92,231,.25)}
.welcome h2{font-size:24px;font-weight:700;letter-spacing:-.5px}
.welcome p{font-size:14px;color:var(--text2);max-width:360px;line-height:1.6}
.suggestions{display:flex;flex-wrap:wrap;gap:8px;justify-content:center;margin-top:4px}
.suggestion{padding:10px 18px;border-radius:24px;background:var(--surface2);border:1px solid var(--border);font-size:13px;color:var(--text2);cursor:pointer;transition:all .2s ease}
.suggestion:hover{border-color:var(--accent);color:var(--accent2);background:var(--accent-glow)}

/* Input */
.input-area{padding:12px 16px;padding-bottom:calc(12px + var(--safe-bottom));background:var(--surface);border-top:1px solid var(--border);flex-shrink:0}
.input-row{display:flex;align-items:flex-end;gap:10px;background:var(--surface2);border:1px solid var(--border);border-radius:24px;padding:6px 6px 6px 18px;transition:border-color .2s,box-shadow .2s}
.input-row:focus-within{border-color:var(--accent);box-shadow:0 0 0 3px rgba(108,92,231,.1)}
.input-row textarea{flex:1;border:none;outline:none;background:transparent;color:var(--text);font-size:14px;font-family:inherit;resize:none;min-height:24px;max-height:120px;line-height:24px;padding:4px 0}
.input-row textarea::placeholder{color:var(--text3)}
.send-btn{width:38px;height:38px;border-radius:50%;border:none;background:var(--accent);color:#fff;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:all .2s;flex-shrink:0}
.send-btn:hover{background:var(--accent2);transform:scale(1.05)}
.send-btn:active{transform:scale(.92)}
.send-btn:disabled{opacity:.3;cursor:not-allowed;transform:none}
.send-btn svg{width:16px;height:16px}
.toast{position:fixed;top:20px;left:50%;transform:translateX(-50%);background:var(--surface2);border:1px solid var(--border);padding:10px 20px;border-radius:12px;font-size:13px;color:var(--text2);opacity:0;transition:opacity .3s;z-index:200;pointer-events:none}
.toast.show{opacity:1}
.overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:40;opacity:0;transition:opacity .3s}
.overlay.show{display:block;opacity:1}
@media(max-width:768px){.sidebar{position:fixed;left:0;top:0;bottom:0;transform:translateX(-100%)}.sidebar.open{transform:translateX(0)}.menu-btn{display:flex}.msg{max-width:92%}}
@media(min-width:769px){.menu-btn{display:none}}
@media(max-width:380px){.header{padding:0 12px;height:52px}.messages{padding:14px 10px;gap:8px}.msg-bubble{padding:10px 14px;font-size:13px}.suggestion{padding:8px 14px;font-size:12px}.input-area{padding:8px 10px;padding-bottom:calc(8px + var(--safe-bottom))}}
</style>
</head>
<body>
<div class="app">
  <aside class="sidebar" id="sidebar">
    <div class="sidebar-header">
      <div class="sidebar-logo">A</div>
      <span class="sidebar-brand">Assistant</span>
    </div>
    <button class="new-chat-btn" onclick="newChat()">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      New Chat
    </button>
    <div class="conv-list" id="convList"></div>
    <div class="sidebar-footer">
      <span>Storage</span>
      <span class="db-badge off" id="dbBadge">Local</span>
    </div>
  </aside>
  <div class="overlay" id="overlay" onclick="closeSidebar()"></div>
  <div class="main">
    <div class="header">
      <button class="menu-btn" id="menuBtn" onclick="toggleSidebar()">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
      </button>
      <div class="header-title" id="headerTitle">New Chat</div>
      <div class="status-pill"><div class="status-dot"></div>Online</div>
    </div>
    <div class="messages" id="messages">
      <div class="welcome" id="welcome">
        <div class="welcome-icon">A</div>
        <h2>Hello</h2>
        <p>How can I help you today? Start a conversation or pick a suggestion below.</p>
        <div class="suggestions">
          <div class="suggestion" onclick="useSuggestion(this)">What can you do?</div>
          <div class="suggestion" onclick="useSuggestion(this)">Summarize a topic</div>
          <div class="suggestion" onclick="useSuggestion(this)">Help me code</div>
          <div class="suggestion" onclick="useSuggestion(this)">Write a story</div>
        </div>
      </div>
      <div class="typing" id="typing">
        <div class="typing-dots"><span></span><span></span><span></span></div>
        <span class="typing-label">Thinking...</span>
      </div>
    </div>
    <div class="input-area">
      <div class="input-row">
        <textarea id="input" rows="1" placeholder="Send a message..." autofocus></textarea>
        <button class="send-btn" id="sendBtn" onclick="sendMessage()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</div>
<div class="toast" id="toast"></div>
<script>
var BT=String.fromCharCode(96);
var currentChatID=null,conversations=[],hasDB=false,eventSource=null,waiting=false;
var $=function(id){return document.getElementById(id)};
var messagesEl=$('messages'),welcomeEl=$('welcome'),typingEl=$('typing'),inputEl=$('input');
var sendBtnEl=$('sendBtn'),convListEl=$('convList'),sidebarEl=$('sidebar'),overlayEl=$('overlay');
var headerTitleEl=$('headerTitle'),dbBadgeEl=$('dbBadge'),toastEl=$('toast');

function toggleSidebar(){sidebarEl.classList.toggle('open');overlayEl.classList.toggle('show')}
function closeSidebar(){sidebarEl.classList.remove('open');overlayEl.classList.remove('show')}
function genID(){return 'web-'+Date.now()+'-'+Math.random().toString(36).slice(2,8)}

function newChat(){
  currentChatID=genID();headerTitleEl.textContent='New Chat';clearMessages();
  welcomeEl.style.display='';connectSSE();closeSidebar();inputEl.focus();
}
function selectChat(id){
  currentChatID=id;
  var c=conversations.find(function(x){return x.id===id});
  if(c)headerTitleEl.textContent=c.title||'Chat';
  clearMessages();welcomeEl.style.display='none';
  loadMessages(id);connectSSE();renderConvList();closeSidebar();inputEl.focus();
}
function deleteChat(e,id){
  e.stopPropagation();
  if(hasDB){fetch('/api/conversations/'+id,{method:'DELETE'})}
  conversations=conversations.filter(function(x){return x.id!==id});
  renderConvList();if(currentChatID===id)newChat();
}
function renderConvList(){
  convListEl.innerHTML='';
  conversations.forEach(function(c){
    var el=document.createElement('div');
    el.className='conv-item'+(c.id===currentChatID?' active':'');
    var t=new Date(c.updated_at||c.created_at);
    var ts=isToday(t)?t.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'}):t.toLocaleDateString([],{month:'short',day:'numeric'});
    el.innerHTML='<span class="conv-title">'+esc(c.title||'Chat')+'</span><span class="conv-time">'+ts+'</span><button class="conv-del" title="Delete">&times;</button>';
    el.onclick=function(){selectChat(c.id)};
    el.querySelector('.conv-del').onclick=function(ev){deleteChat(ev,c.id)};
    convListEl.appendChild(el);
  });
}
function isToday(d){var n=new Date();return d.getDate()===n.getDate()&&d.getMonth()===n.getMonth()&&d.getFullYear()===n.getFullYear()}

function loadConversations(){
  fetch('/api/conversations').then(function(r){return r.json()}).then(function(d){
    hasDB=!!d.db;dbBadgeEl.textContent=hasDB?'PostgreSQL':'Local';
    dbBadgeEl.className='db-badge '+(hasDB?'on':'off');
    if(d.conversations&&d.conversations.length){conversations=d.conversations;renderConvList()}
  }).catch(function(e){console.error('Load conversations',e)});
}
function loadMessages(chatID){
  if(!hasDB)return;
  fetch('/api/conversations/'+chatID+'/messages').then(function(r){return r.json()}).then(function(d){
    if(d.messages){d.messages.forEach(function(m){addMessage(m.content,m.role==='user'?'user':'bot',false)})}
  }).catch(function(e){console.error('Load messages',e)});
}
function connectSSE(){
  if(eventSource)eventSource.close();if(!currentChatID)return;
  eventSource=new EventSource('/api/stream?chat_id='+encodeURIComponent(currentChatID));
  eventSource.onmessage=function(e){
    try{var data=JSON.parse(e.data);
      if(data.type==='message'){hideTyping();addMessage(data.content,'bot',true);waiting=false;updateSendBtn()}
    }catch(err){console.error('SSE parse',err)}
  };
  eventSource.onerror=function(){setTimeout(function(){if(currentChatID)connectSSE()},3000)};
}
function clearMessages(){var ms=messagesEl.querySelectorAll('.msg');ms.forEach(function(m){m.remove()})}
function addMessage(text,role,animate){
  welcomeEl.style.display='none';
  var wrap=document.createElement('div');wrap.className='msg msg--'+role;
  if(!animate)wrap.style.animation='none';
  var avatar=document.createElement('div');avatar.className='msg-avatar';avatar.textContent=role==='user'?'U':'P';
  var body=document.createElement('div');body.className='msg-body';
  var bubble=document.createElement('div');bubble.className='msg-bubble';
  if(role==='bot'){bubble.innerHTML=renderMarkdown(text)}else{bubble.textContent=text}
  var tm=document.createElement('div');tm.className='msg-time';
  tm.textContent=new Date().toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});
  body.appendChild(bubble);body.appendChild(tm);wrap.appendChild(avatar);wrap.appendChild(body);
  messagesEl.insertBefore(wrap,typingEl);scrollBottom();
}
function renderMarkdown(t){
  var h=t.replace(/` + "```" + `(\\w*)\\n([\\s\\S]*?)` + "```" + `/g,function(_,lang,code){
    return '<pre><code>'+esc(code.trim())+'</code><button class="copy-btn" onclick="copyCode(this)">Copy</button></pre>';
  });
  h=h.replace(/` + "`" + `([^` + "`" + `]+)` + "`" + `/g,'<code>$1</code>');
  h=h.replace(/\*\*(.+?)\*\*/g,'<strong>$1</strong>');
  h=h.replace(/\*(.+?)\*/g,'<em>$1</em>');
  h=h.replace(/\n\n/g,'</p><p>');h=h.replace(/\n/g,'<br>');
  h='<p>'+h+'</p>';h=h.replace(/<p><\/p>/g,'');return h;
}
function copyCode(btn){var code=btn.parentElement.querySelector('code').textContent;navigator.clipboard.writeText(code).then(function(){btn.textContent='Copied!';setTimeout(function(){btn.textContent='Copy'},1500)})}
function esc(s){var d=document.createElement('div');d.textContent=s;return d.innerHTML}
function showTyping(){typingEl.classList.add('active');scrollBottom()}
function hideTyping(){typingEl.classList.remove('active')}
function scrollBottom(){requestAnimationFrame(function(){messagesEl.scrollTop=messagesEl.scrollHeight})}
function updateSendBtn(){sendBtnEl.disabled=waiting||!inputEl.value.trim()}

function sendMessage(){
  var text=inputEl.value.trim();if(!text||waiting)return;
  if(!currentChatID)currentChatID=genID();
  inputEl.value='';autoResize();addMessage(text,'user',true);
  waiting=true;updateSendBtn();showTyping();
  var existing=conversations.find(function(c){return c.id===currentChatID});
  if(!existing){
    var title=text.length>60?text.slice(0,60)+'...':text;
    conversations.unshift({id:currentChatID,title:title,updated_at:new Date().toISOString()});
    headerTitleEl.textContent=title;
  }else{existing.updated_at=new Date().toISOString()}
  renderConvList();
  fetch('/api/chat',{method:'POST',headers:{'Content-Type':'application/json'},
    body:JSON.stringify({message:text,chat_id:currentChatID})
  }).then(function(r){if(!r.ok)return r.text().then(function(t){throw new Error(t)})})
  .catch(function(err){hideTyping();waiting=false;updateSendBtn();showToast('Failed to send. Try again.');console.error('Send',err)});
}
function useSuggestion(el){inputEl.value=el.textContent;sendMessage()}
function showToast(msg){toastEl.textContent=msg;toastEl.classList.add('show');setTimeout(function(){toastEl.classList.remove('show')},3000)}
function autoResize(){inputEl.style.height='auto';inputEl.style.height=Math.min(inputEl.scrollHeight,120)+'px'}
inputEl.addEventListener('input',function(){autoResize();updateSendBtn()});
inputEl.addEventListener('keydown',function(e){if(e.key==='Enter'&&!e.shiftKey){e.preventDefault();sendMessage()}});
newChat();loadConversations();updateSendBtn();
</script>
</body>
</html>`
