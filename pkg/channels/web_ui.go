package channels

const webUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
<title>PicoClaw</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0a0a0f;
  --surface:#12121a;
  --surface2:#1a1a26;
  --border:#2a2a3a;
  --text:#e4e4ef;
  --text2:#8888a0;
  --accent:#6c5ce7;
  --accent2:#a29bfe;
  --user-bg:#6c5ce7;
  --bot-bg:#1e1e2e;
  --radius:16px;
  --safe-bottom:env(safe-area-inset-bottom,0px);
}
html,body{
  height:100%;width:100%;
  font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;
  background:var(--bg);color:var(--text);
  overflow:hidden;
  -webkit-font-smoothing:antialiased;
}
.app{display:flex;flex-direction:column;height:100%;max-width:900px;margin:0 auto}

/* Header */
.header{
  display:flex;align-items:center;gap:12px;
  padding:16px 20px;
  background:var(--surface);
  border-bottom:1px solid var(--border);
  flex-shrink:0;
  backdrop-filter:blur(20px);
  -webkit-backdrop-filter:blur(20px);
}
.header-logo{
  width:40px;height:40px;border-radius:12px;
  background:linear-gradient(135deg,var(--accent),var(--accent2));
  display:flex;align-items:center;justify-content:center;
  font-size:20px;font-weight:700;color:#fff;
  box-shadow:0 4px 12px rgba(108,92,231,.3);
}
.header-info h1{font-size:17px;font-weight:600;letter-spacing:-.3px}
.header-info p{font-size:12px;color:var(--text2);margin-top:1px}
.status-dot{
  width:8px;height:8px;border-radius:50%;
  background:#2ed573;
  display:inline-block;margin-right:4px;
  animation:pulse 2s ease infinite;
}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}

/* Messages area */
.messages{
  flex:1;overflow-y:auto;overflow-x:hidden;
  padding:20px 16px;
  display:flex;flex-direction:column;gap:8px;
  scroll-behavior:smooth;
  -webkit-overflow-scrolling:touch;
}
.messages::-webkit-scrollbar{width:4px}
.messages::-webkit-scrollbar-track{background:transparent}
.messages::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}

.msg{
  display:flex;gap:8px;
  max-width:85%;
  animation:msgIn .3s ease;
}
@keyframes msgIn{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}
.msg--user{align-self:flex-end;flex-direction:row-reverse}
.msg--bot{align-self:flex-start}

.msg-avatar{
  width:32px;height:32px;border-radius:10px;flex-shrink:0;
  display:flex;align-items:center;justify-content:center;
  font-size:14px;
}
.msg--bot .msg-avatar{background:var(--surface2);color:var(--accent2)}
.msg--user .msg-avatar{background:var(--user-bg);color:#fff}

.msg-bubble{
  padding:12px 16px;
  border-radius:var(--radius);
  font-size:15px;line-height:1.5;
  word-wrap:break-word;
  white-space:pre-wrap;
}
.msg--user .msg-bubble{
  background:var(--user-bg);color:#fff;
  border-bottom-right-radius:4px;
}
.msg--bot .msg-bubble{
  background:var(--bot-bg);color:var(--text);
  border:1px solid var(--border);
  border-bottom-left-radius:4px;
}

/* Typing indicator */
.typing{display:none;align-self:flex-start;padding:0 16px}
.typing.active{display:flex;align-items:center;gap:8px}
.typing-dots{display:flex;gap:4px}
.typing-dots span{
  width:6px;height:6px;border-radius:50%;
  background:var(--text2);
  animation:typingBounce 1.4s ease infinite;
}
.typing-dots span:nth-child(2){animation-delay:.2s}
.typing-dots span:nth-child(3){animation-delay:.4s}
@keyframes typingBounce{0%,60%,100%{transform:translateY(0)}30%{transform:translateY(-6px)}}
.typing-label{font-size:12px;color:var(--text2)}

/* Welcome */
.welcome{
  flex:1;display:flex;flex-direction:column;
  align-items:center;justify-content:center;
  gap:16px;padding:40px 20px;text-align:center;
}
.welcome-icon{
  width:72px;height:72px;border-radius:20px;
  background:linear-gradient(135deg,var(--accent),var(--accent2));
  display:flex;align-items:center;justify-content:center;
  font-size:36px;color:#fff;
  box-shadow:0 8px 32px rgba(108,92,231,.25);
}
.welcome h2{font-size:22px;font-weight:700;letter-spacing:-.5px}
.welcome p{font-size:14px;color:var(--text2);max-width:320px;line-height:1.5}
.suggestions{display:flex;flex-wrap:wrap;gap:8px;justify-content:center;margin-top:8px}
.suggestion{
  padding:8px 16px;border-radius:20px;
  background:var(--surface2);border:1px solid var(--border);
  font-size:13px;color:var(--text2);cursor:pointer;
  transition:all .2s ease;
}
.suggestion:hover{border-color:var(--accent);color:var(--accent2);background:rgba(108,92,231,.08)}

/* Input area */
.input-area{
  padding:12px 16px;
  padding-bottom:calc(12px + var(--safe-bottom));
  background:var(--surface);
  border-top:1px solid var(--border);
  flex-shrink:0;
}
.input-row{
  display:flex;align-items:flex-end;gap:10px;
  background:var(--surface2);
  border:1px solid var(--border);
  border-radius:24px;
  padding:6px 6px 6px 18px;
  transition:border-color .2s ease;
}
.input-row:focus-within{border-color:var(--accent)}
.input-row textarea{
  flex:1;border:none;outline:none;
  background:transparent;color:var(--text);
  font-size:15px;font-family:inherit;
  resize:none;min-height:24px;max-height:120px;
  line-height:24px;padding:4px 0;
}
.input-row textarea::placeholder{color:var(--text2)}
.send-btn{
  width:40px;height:40px;border-radius:50%;border:none;
  background:var(--accent);color:#fff;cursor:pointer;
  display:flex;align-items:center;justify-content:center;
  transition:all .2s ease;flex-shrink:0;
}
.send-btn:hover{background:var(--accent2);transform:scale(1.05)}
.send-btn:active{transform:scale(.95)}
.send-btn:disabled{opacity:.4;cursor:not-allowed;transform:none}
.send-btn svg{width:18px;height:18px}

/* Toast */
.toast{
  position:fixed;top:20px;left:50%;transform:translateX(-50%);
  background:var(--surface2);border:1px solid var(--border);
  padding:10px 20px;border-radius:12px;
  font-size:13px;color:var(--text2);
  opacity:0;transition:opacity .3s ease;
  z-index:100;pointer-events:none;
}
.toast.show{opacity:1}

/* Desktop */
@media(min-width:768px){
  .app{border-left:1px solid var(--border);border-right:1px solid var(--border)}
  .msg{max-width:70%}
}
/* Small phones */
@media(max-width:380px){
  .header{padding:12px 14px}
  .messages{padding:14px 10px}
  .msg-bubble{padding:10px 14px;font-size:14px}
  .suggestions{gap:6px}
  .suggestion{padding:6px 12px;font-size:12px}
}
</style>
</head>
<body>
<div class="app">
  <div class="header">
    <div class="header-logo">P</div>
    <div class="header-info">
      <h1>PicoClaw</h1>
      <p><span class="status-dot"></span>Online</p>
    </div>
  </div>

  <div class="messages" id="messages">
    <div class="welcome" id="welcome">
      <div class="welcome-icon">P</div>
      <h2>PicoClaw Agent</h2>
      <p>Your ultra-lightweight personal AI assistant. Ask me anything or try a suggestion below.</p>
      <div class="suggestions">
        <div class="suggestion" onclick="sendSuggestion(this)">What can you do?</div>
        <div class="suggestion" onclick="sendSuggestion(this)">Tell me a joke</div>
        <div class="suggestion" onclick="sendSuggestion(this)">Help me code</div>
      </div>
    </div>
    <div class="typing" id="typing">
      <div class="typing-dots"><span></span><span></span><span></span></div>
      <span class="typing-label">PicoClaw is thinking...</span>
    </div>
  </div>

  <div class="input-area">
    <div class="input-row">
      <textarea id="input" rows="1" placeholder="Message PicoClaw..." autofocus></textarea>
      <button class="send-btn" id="sendBtn" onclick="sendMessage()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>
        </svg>
      </button>
    </div>
  </div>
</div>
<div class="toast" id="toast"></div>

<script>
const chatID='web-'+Date.now()+'-'+Math.random().toString(36).slice(2,8);
const messagesEl=document.getElementById('messages');
const welcomeEl=document.getElementById('welcome');
const typingEl=document.getElementById('typing');
const inputEl=document.getElementById('input');
const sendBtnEl=document.getElementById('sendBtn');
const toastEl=document.getElementById('toast');
let eventSource=null;
let waiting=false;

function connectSSE(){
  if(eventSource)eventSource.close();
  eventSource=new EventSource('/api/stream?chat_id='+encodeURIComponent(chatID));
  eventSource.onmessage=function(e){
    try{
      const data=JSON.parse(e.data);
      if(data.type==='message'){
        hideTyping();
        addMessage(data.content,'bot');
        waiting=false;
        updateSendBtn();
      }
    }catch(err){console.error('SSE parse error',err)}
  };
  eventSource.onerror=function(){
    setTimeout(connectSSE,3000);
  };
}

function addMessage(text,role){
  welcomeEl.style.display='none';
  const wrap=document.createElement('div');
  wrap.className='msg msg--'+role;
  const avatar=document.createElement('div');
  avatar.className='msg-avatar';
  avatar.textContent=role==='user'?'U':'P';
  const bubble=document.createElement('div');
  bubble.className='msg-bubble';
  bubble.textContent=text;
  wrap.appendChild(avatar);
  wrap.appendChild(bubble);
  messagesEl.insertBefore(wrap,typingEl);
  scrollBottom();
}

function showTyping(){typingEl.classList.add('active');scrollBottom()}
function hideTyping(){typingEl.classList.remove('active')}

function scrollBottom(){
  requestAnimationFrame(()=>{messagesEl.scrollTop=messagesEl.scrollHeight});
}

function updateSendBtn(){
  sendBtnEl.disabled=waiting||!inputEl.value.trim();
}

async function sendMessage(){
  const text=inputEl.value.trim();
  if(!text||waiting)return;
  inputEl.value='';
  autoResize();
  addMessage(text,'user');
  waiting=true;
  updateSendBtn();
  showTyping();
  try{
    const resp=await fetch('/api/chat',{
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body:JSON.stringify({message:text,chat_id:chatID})
    });
    if(!resp.ok){
      const errText=await resp.text();
      throw new Error(errText);
    }
  }catch(err){
    hideTyping();
    waiting=false;
    updateSendBtn();
    showToast('Failed to send message. Please try again.');
    console.error('Send error',err);
  }
}

function sendSuggestion(el){
  inputEl.value=el.textContent;
  sendMessage();
}

function showToast(msg){
  toastEl.textContent=msg;
  toastEl.classList.add('show');
  setTimeout(()=>toastEl.classList.remove('show'),3000);
}

function autoResize(){
  inputEl.style.height='auto';
  inputEl.style.height=Math.min(inputEl.scrollHeight,120)+'px';
}

inputEl.addEventListener('input',function(){autoResize();updateSendBtn()});
inputEl.addEventListener('keydown',function(e){
  if(e.key==='Enter'&&!e.shiftKey){e.preventDefault();sendMessage()}
});

updateSendBtn();
connectSSE();
</script>
</body>
</html>`
