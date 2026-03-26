'use strict';

// ============================================================
// Batch list — must match server Batches slice
// ============================================================
const BATCHES = [10,25,50,100,250,500,1000,2500,5000,10000,25000,50000,100000,250000,500000,1000000];

// ============================================================
// State
// ============================================================
let numDoors    = 3;
let simDelay    = 80;
let maxTrials   = 1_000_000;
let useParallel = false;
let eventSource = null;
let isRunning   = false;
let simResults  = [];   // for CSV download

// Demo state
let demoTimers  = [];

// Play mode state
let playMode    = false;
let playState   = 'idle'; // idle | picked | revealed | done
let playCarPos  = -1;
let playPicked  = -1;
let playMonty   = -1;
let playSwitchDoor = -1;
let playTimer   = null;
let playStats   = { sw: [0,0], st: [0,0] }; // [wins, total]

// ============================================================
// DOM refs
// ============================================================
const $ = id => document.getElementById(id);

const doorsRange    = $('doors-range');
const doorsDisplay  = $('doors-display');
const maxSelect     = $('max-select');
const speedSelect   = $('speed-select');
const parallelCheck = $('parallel-check');
const startBtn      = $('start-btn');
const resetBtn      = $('reset-btn');
const downloadBtn   = $('download-btn');
const progressWrap  = $('progress-wrap');
const progressBar   = $('progress-bar');
const trialLabel    = $('trial-label');

const stayPctEl   = $('stay-pct');
const stayCiEl    = $('stay-ci');
const stayTheoEl  = $('stay-theory');
const stayDeltaEl = $('stay-delta');
const swPctEl     = $('switch-pct');
const swCiEl      = $('switch-ci');
const swTheoEl    = $('switch-theory');
const swDeltaEl   = $('switch-delta');
const gainEl      = $('gain-value');

const mathBtn     = $('math-btn');
const mathBody    = $('math-body');
const mathContent = $('math-content');

const themeBtn    = $('theme-btn');

// ============================================================
// Dark mode
// ============================================================
function isDark() {
  return document.documentElement.getAttribute('data-theme') === 'dark';
}

function setTheme(dark) {
  document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
  localStorage.setItem('theme', dark ? 'dark' : 'light');
  updateChartTheme(dark);
}

function updateChartTheme(dark) {
  const gridColor = dark ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)';
  const textColor = dark ? '#9e9087' : '#5c4a3e';
  ['x','y'].forEach(ax => {
    chart.options.scales[ax].grid  = { color: gridColor };
    chart.options.scales[ax].ticks = { ...chart.options.scales[ax].ticks, color: textColor };
    chart.options.scales[ax].title.color = textColor;
  });
  chart.options.plugins.legend.labels.color = textColor;
  chart.update('none');
}

themeBtn.addEventListener('click', () => setTheme(!isDark()));

// Theme init deferred to after chart creation (see Init section at bottom)

// ============================================================
// Chart
// ============================================================
const chart = new Chart($('chart').getContext('2d'), {
  type: 'line',
  data: {
    labels: [],
    datasets: [
      { label: 'Switch (simulated)',  data: [], borderColor: '#f59e0b', backgroundColor: 'rgba(245,158,11,0.08)', borderWidth: 3, pointRadius: 5, pointHoverRadius: 8, tension: 0.3, fill: false },
      { label: 'Stay (simulated)',    data: [], borderColor: '#60a5fa', backgroundColor: 'rgba(96,165,250,0.08)',  borderWidth: 3, pointRadius: 5, pointHoverRadius: 8, tension: 0.3, fill: false },
      { label: 'Theoretical Switch',  data: [], borderColor: 'rgba(245,158,11,0.40)', borderDash: [8,5], borderWidth: 2, pointRadius: 0, fill: false },
      { label: 'Theoretical Stay',    data: [], borderColor: 'rgba(96,165,250,0.40)',  borderDash: [8,5], borderWidth: 2, pointRadius: 0, fill: false },
    ],
  },
  options: {
    responsive: true, maintainAspectRatio: false,
    animation: { duration: 250 },
    interaction: { mode: 'index', intersect: false },
    plugins: {
      legend: { position: 'top', labels: { usePointStyle: true, padding: 22, font: { size: 13 } } },
      tooltip: { callbacks: { label: ctx => `${ctx.dataset.label}: ${ctx.parsed.y.toFixed(2)}%` } },
    },
    scales: {
      x: { title: { display: true, text: 'Number of simulations', font: { size: 13 } }, ticks: { maxTicksLimit: 8, maxRotation: 30 } },
      y: { min: 0, max: 100, title: { display: true, text: 'Win rate (%)', font: { size: 13 } }, ticks: { callback: v => v + '%' }, grid: { color: 'rgba(0,0,0,0.04)' } },
    },
  },
});

function resetChart() {
  chart.data.labels = [];
  chart.data.datasets.forEach(ds => (ds.data = []));
  chart.update('none');
}

function pushToChart(d) {
  chart.data.labels.push(fmtCount(d.trial_count));
  chart.data.datasets[0].data.push(+d.switch_pct.toFixed(3));
  chart.data.datasets[1].data.push(+d.stay_pct.toFixed(3));
  chart.data.datasets[2].data.push(+d.theo_switch.toFixed(3));
  chart.data.datasets[3].data.push(+d.theo_stay.toFixed(3));
  chart.data.datasets[2].label = `Theoretical Switch (${d.theo_switch.toFixed(2)}%)`;
  chart.data.datasets[3].label = `Theoretical Stay (${d.theo_stay.toFixed(2)}%)`;
  chart.update();
}

// ============================================================
// Door helpers (shared by demo and play mode)
// ============================================================
function buildDoors(n) {
  const w = Math.max(60, Math.min(88, Math.floor(300 / n)));
  const h = Math.round(w * 1.55);
  const container = $('doors-container');
  container.innerHTML = '';
  container.style.setProperty('--door-w', w + 'px');
  container.style.setProperty('--door-h', h + 'px');

  for (let i = 0; i < n; i++) {
    const wrap = document.createElement('div');
    wrap.className = 'door-wrap';
    wrap.dataset.i = i;
    wrap.innerHTML = `
      <div class="door-inner">
        <div class="door-front">
          <span class="door-star">&#x2605;</span>
          <span class="door-number">${i + 1}</span>
        </div>
        <div class="door-back"><span class="door-emoji"></span></div>
      </div>`;
    wrap.addEventListener('click', () => onDoorClick(i));
    container.appendChild(wrap);
  }
}

function doorEl(i) { return document.querySelector(`.door-wrap[data-i="${i}"]`); }

function flipDoor(i, emoji, state) {
  const el = doorEl(i);
  if (!el) return;
  el.querySelector('.door-back .door-emoji').textContent = emoji;
  el.classList.remove('revealed', 'win', 'lose');
  if (state) el.classList.add(state);
  el.classList.add('flipped');
}

function pickDoor(i) {
  document.querySelectorAll('.door-wrap').forEach(el => el.classList.remove('picked'));
  const el = doorEl(i);
  if (el) el.classList.add('picked');
}

function resetDoors() {
  document.querySelectorAll('.door-wrap').forEach(el => {
    el.classList.remove('flipped', 'picked', 'revealed', 'win', 'lose', 'clickable');
    el.querySelector('.door-back .door-emoji').textContent = '?';
  });
}

function setDoorsClickable(on) {
  document.querySelectorAll('.door-wrap').forEach(el =>
    el.classList.toggle('clickable', on)
  );
}

// ============================================================
// Demo mode (auto-loop)
// ============================================================
function stopDemo() { demoTimers.forEach(clearTimeout); demoTimers = []; }

function runDemo() {
  if (playMode) return;
  stopDemo();
  resetDoors();
  $('demo-result').className = 'demo-result hidden';

  const n          = numDoors;
  const carPos     = rnd(n);
  const pick       = rnd(n);
  const hint       = $('demo-hint');

  let monty = -1;
  for (let i = 0; i < n; i++) { if (i !== pick && i !== carPos) { monty = i; break; } }
  let sw = -1;
  for (let i = 0; i < n; i++) { if (i !== pick && i !== monty) { sw = i; break; } }

  const t1 = setTimeout(() => {
    hint.textContent = `You pick door ${pick + 1}\u2026`;
    pickDoor(pick);
  }, 300);

  const t2 = setTimeout(() => {
    hint.textContent = `Host opens door ${monty + 1}: it\u2019s a goat!`;
    flipDoor(monty, '\uD83D\uDC10', 'revealed');
  }, 1200);

  const t3 = setTimeout(() => {
    hint.textContent = `Switching to door ${sw + 1}\u2026`;
    const wins = sw === carPos;
    flipDoor(sw,  wins ? '\uD83D\uDE97' : '\uD83D\uDC10', wins ? 'win' : 'lose');
    flipDoor(pick, pick === carPos ? '\uD83D\uDE97' : '\uD83D\uDC10', pick === carPos ? 'win' : 'lose');

    const res = $('demo-result');
    res.textContent = wins
      ? `\uD83C\uDF89 Switch wins! Door ${sw + 1} had the car!`
      : `\uD83D\uDE14 This time staying would have won.`;
    res.className = `demo-result fade-up ${wins ? 'switch-win' : 'stay-win'}`;
  }, 2200);

  const t4 = setTimeout(() => runDemo(), 5000);
  demoTimers = [t1, t2, t3, t4];
}

// ============================================================
// Play mode
// ============================================================
function setMode(play) {
  playMode = play;
  $('btn-watch').classList.toggle('active', !play);
  $('btn-play').classList.toggle('active', play);

  if (play) {
    stopDemo();
    resetDoors();
    $('demo-title').textContent = '\uD83C\uDFAE Play it Yourself!';
    $('play-stats').classList.remove('hidden');
    $('demo-hint').textContent = 'Click a door to pick it!';
    $('demo-result').className = 'demo-result hidden';
    $('play-choice').classList.add('hidden');
    playState = 'idle';
    setDoorsClickable(true);
  } else {
    if (playTimer) clearTimeout(playTimer);
    playState = 'idle';
    $('demo-title').textContent = '\uD83C\uDF9C Live Demo';
    $('play-stats').classList.add('hidden');
    $('play-choice').classList.add('hidden');
    setDoorsClickable(false);
    setTimeout(() => runDemo(), 200);
  }
}

function onDoorClick(i) {
  if (!playMode || playState !== 'idle') return;

  playCarPos = rnd(numDoors);
  playPicked = i;
  playState  = 'picked';

  pickDoor(i);
  $('demo-hint').textContent = `You picked door ${i + 1}. Host is thinking\u2026`;

  playMonty = -1;
  for (let d = 0; d < numDoors; d++) {
    if (d !== playPicked && d !== playCarPos) { playMonty = d; break; }
  }
  playSwitchDoor = -1;
  for (let d = 0; d < numDoors; d++) {
    if (d !== playPicked && d !== playMonty) { playSwitchDoor = d; break; }
  }

  playTimer = setTimeout(() => {
    flipDoor(playMonty, '\uD83D\uDC10', 'revealed');
    $('demo-hint').textContent =
      `Host opens door ${playMonty + 1}: a goat! Switch to ${playSwitchDoor + 1} or stay?`;
    $('play-choice').classList.remove('hidden');
    playState = 'revealed';
  }, 700);
}

function onPlayChoice(choice) {
  if (!playMode || playState !== 'revealed') return;
  playState = 'done';
  $('play-choice').classList.add('hidden');

  const finalDoor = choice === 'switch' ? playSwitchDoor : playPicked;
  const wins      = finalDoor === playCarPos;

  if (choice === 'switch') { playStats.sw[1]++; if (wins) playStats.sw[0]++; }
  else                     { playStats.st[1]++; if (wins) playStats.st[0]++; }

  flipDoor(playSwitchDoor, playSwitchDoor === playCarPos ? '\uD83D\uDE97' : '\uD83D\uDC10',
    playSwitchDoor === playCarPos ? 'win' : 'lose');
  flipDoor(playPicked, playPicked === playCarPos ? '\uD83D\uDE97' : '\uD83D\uDC10',
    playPicked === playCarPos ? 'win' : 'lose');

  const res = $('demo-result');
  res.textContent = wins
    ? `\uD83C\uDF89 You win! Door ${finalDoor + 1} had the car!`
    : `\uD83D\uDE14 You lose. The car was behind door ${playCarPos + 1}.`;
  res.className = `demo-result fade-up ${wins ? 'switch-win' : 'stay-win'}`;
  $('demo-hint').textContent = `${choice === 'switch' ? 'You switched' : 'You stayed'}: ${wins ? 'WIN!' : 'LOSE.'}`;

  updatePlayStats();

  playTimer = setTimeout(() => {
    if (!playMode) return;
    resetDoors();
    setDoorsClickable(true);
    $('demo-result').className = 'demo-result hidden';
    $('play-choice').classList.add('hidden');
    $('demo-hint').textContent = 'Click a door to pick it!';
    playState = 'idle';
  }, 2500);
}

function updatePlayStats() {
  const fmt = ([w, t]) => t > 0 ? `${w}/${t} (${Math.round(w/t*100)}%)` : '0/0';
  $('play-games').textContent       = playStats.sw[1] + playStats.st[1];
  $('play-switch-stat').textContent = fmt(playStats.sw);
  $('play-stay-stat').textContent   = fmt(playStats.st);
}

// ============================================================
// Simulation
// ============================================================
function activeBatchCount(max) {
  return BATCHES.filter(n => n <= max).length || 1;
}

function startSim() {
  if (isRunning) return;
  stopDemo();
  resetChart();
  simResults  = [];
  isRunning   = true;

  startBtn.disabled = true;
  startBtn.innerHTML = '&#x23F3; Running&hellip;';
  downloadBtn.classList.add('hidden');
  progressWrap.classList.remove('hidden');
  progressBar.style.width = '0%';

  const url = `/simulate?doors=${numDoors}&delay=${simDelay}&max=${maxTrials}&parallel=${useParallel ? 1 : 0}`;
  eventSource = new EventSource(url);

  eventSource.onmessage = e => {
    const data = JSON.parse(e.data);
    simResults.push(data);
    updateStats(data);
    pushToChart(data);
    updateProgress(data);
    if (data.done) finishSim();
  };

  eventSource.onerror = () => finishSim();
}

function finishSim() {
  if (eventSource) { eventSource.close(); eventSource = null; }
  isRunning = false;
  startBtn.disabled = false;
  startBtn.innerHTML = '&#x25B6; Run Again';
  progressBar.style.width = '100%';
  if (simResults.length > 0) downloadBtn.classList.remove('hidden');
  if (!playMode) setTimeout(() => runDemo(), 1200);
}

function resetSim() {
  if (eventSource) { eventSource.close(); eventSource = null; }
  isRunning  = false;
  simResults = [];
  startBtn.disabled = false;
  startBtn.innerHTML = '&#x25B6; Start';
  progressWrap.classList.add('hidden');
  downloadBtn.classList.add('hidden');
  progressBar.style.width = '0%';
  trialLabel.textContent = '0 simulations';
  resetChart();
  clearStats();
}

function updateProgress(d) {
  const total = activeBatchCount(maxTrials);
  const done  = simResults.length;
  progressBar.style.width = Math.min(100, done / total * 100) + '%';
  trialLabel.textContent  = fmtCount(d.trial_count) + ' simulations';
}

// ============================================================
// Stats display
// ============================================================
function updateStats(d) {
  const sd  = Math.abs(d.stay_pct   - d.theo_stay).toFixed(2);
  const swd = Math.abs(d.switch_pct - d.theo_switch).toFixed(2);

  stayPctEl.textContent   = d.stay_pct.toFixed(2) + '%';
  stayCiEl.textContent    = `95% CI: \xB1${d.stay_ci95.toFixed(2)}%`;
  stayDeltaEl.textContent = `\u0394 theory: ${sd}%`;
  stayDeltaEl.style.color = parseFloat(sd) < 1 ? 'var(--ok)' : 'var(--warn)';

  swPctEl.textContent     = d.switch_pct.toFixed(2) + '%';
  swCiEl.textContent      = `95% CI: \xB1${d.switch_ci95.toFixed(2)}%`;
  swDeltaEl.textContent   = `\u0394 theory: ${swd}%`;
  swDeltaEl.style.color   = parseFloat(swd) < 1 ? 'var(--ok)' : 'var(--warn)';

  gainEl.textContent      = '+' + (d.switch_pct - d.stay_pct).toFixed(2) + '%';
  stayTheoEl.textContent  = d.theo_stay.toFixed(2) + '%';
  swTheoEl.textContent    = d.theo_switch.toFixed(2) + '%';

  document.querySelectorAll('.stat-card').forEach(c => {
    c.classList.remove('updated'); void c.offsetWidth; c.classList.add('updated');
  });
}

function clearStats() {
  [stayPctEl, stayCiEl, stayDeltaEl, swPctEl, swCiEl, swDeltaEl].forEach(el => (el.textContent = '\u2014'));
  gainEl.textContent = '\u2014';
  stayDeltaEl.style.color = swDeltaEl.style.color = '';
}

// ============================================================
// CSV download
// ============================================================
function downloadCSV() {
  if (!simResults.length) return;
  const headers = ['trials','stay_pct','switch_pct','stay_ci95','switch_ci95','theo_stay','theo_switch'];
  const rows    = simResults.map(d => [
    d.trial_count,
    d.stay_pct.toFixed(4), d.switch_pct.toFixed(4),
    d.stay_ci95.toFixed(4), d.switch_ci95.toFixed(4),
    d.theo_stay.toFixed(4), d.theo_switch.toFixed(4),
  ]);
  const csv  = [headers, ...rows].map(r => r.join(',')).join('\n');
  const blob = new Blob([csv], { type: 'text/csv' });
  const url  = URL.createObjectURL(blob);
  const a    = document.createElement('a');
  a.href = url;
  a.download = `monty-hall-${numDoors}doors-${fmtCount(maxTrials)}trials.csv`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

// ============================================================
// Math explanation
// ============================================================
function renderMath(n) {
  const ts  = (100 / n).toFixed(4);
  const tsw = (100 * (n - 1) / n).toFixed(4);

  const bayes = n === 3 ? `
    <h3>Bayesian derivation (N = 3)</h3>
    <p>You pick door 1. Host opens door 3 (goat). What is P(car = door 2)?</p>
    <div class="formula">P(car=2 | host opens 3) = P(host opens 3 | car=2) &times; P(car=2) / P(host opens 3)

  P(car=2)               = 1/3
  P(host opens 3 | car=2) = 1     (only option for host)
  P(host opens 3 | car=1) = 1/2   (host can open door 2 or 3)
  P(host opens 3 | car=3) = 0     (cannot reveal the car)

  P(host opens 3) = (1/2)(1/3) + (1)(1/3) + (0)(1/3) = 1/2

  P(car=2 | host opens 3) = (1 &times; 1/3) / (1/2) = <strong>2/3 = 66.67%</strong>
  P(car=1 | host opens 3) = (1/2 &times; 1/3) / (1/2) = <strong>1/3 = 33.33%</strong></div>` : '';

  const tableRows = buildMathTable(n);

  mathContent.innerHTML = `
    <h3>Core result</h3>
    <div class="formula">P(win | stay)   = 1/${n} = ${ts}%
P(win | switch) = ${n-1}/${n} = ${tsw}%</div>
    <h3>Intuition</h3>
    <p>You win by switching <em>if and only if</em> your initial pick was wrong.
    The probability of picking wrong first is <strong>${n-1}/${n}</strong>.
    After the host reveals all other goats, the remaining door concentrates that full probability mass.</p>
    ${bayes}
    <h3>Generalization to N doors</h3>
    <p>With ${n} doors, switching is <strong>${n-1}&times; more likely to win</strong> than staying.</p>
    <div class="table-wrap"><table class="math-table">
      <thead><tr><th>N doors</th><th class="c-blue">P(win | stay)</th><th class="c-gold">P(win | switch)</th><th>Advantage</th></tr></thead>
      <tbody>${tableRows}</tbody>
    </table></div>`;
}

function buildMathTable(current) {
  const set = new Set([3,4,5,10,100,current]);
  return [...set].sort((a,b)=>a-b).map(d => {
    const active = d === current ? ' class="active-row"' : '';
    const mark   = d === current ? ' &#x25C4;' : '';
    return `<tr${active}><td>${d}${mark}</td><td class="c-blue">${(100/d).toFixed(2)}%</td><td class="c-gold">${(100*(d-1)/d).toFixed(2)}%</td><td>${d-1}&times; better</td></tr>`;
  }).join('');
}

// ============================================================
// Utilities
// ============================================================
function rnd(n) { return Math.floor(Math.random() * n); }
function fmtCount(n) { return n.toLocaleString('en-US'); }

function updateTheoDisplay() {
  stayTheoEl.textContent = (100 / numDoors).toFixed(2) + '%';
  swTheoEl.textContent   = (100 * (numDoors - 1) / numDoors).toFixed(2) + '%';
}

// ============================================================
// Event listeners
// ============================================================
doorsRange.addEventListener('input', () => {
  numDoors = parseInt(doorsRange.value, 10);
  doorsDisplay.textContent = numDoors;
  buildDoors(numDoors);
  if (playMode) setDoorsClickable(true);
  updateTheoDisplay();
  renderMath(numDoors);
  if (isRunning) resetSim();
  if (!playMode) { stopDemo(); setTimeout(() => runDemo(), 200); }
});

maxSelect.addEventListener('change',     () => { maxTrials   = parseInt(maxSelect.value, 10); });
speedSelect.addEventListener('change',   () => { simDelay    = parseInt(speedSelect.value, 10); });
parallelCheck.addEventListener('change', () => { useParallel = parallelCheck.checked; });

startBtn.addEventListener('click',    startSim);
resetBtn.addEventListener('click',    () => { resetSim(); stopDemo(); buildDoors(numDoors); if (!playMode) setTimeout(() => runDemo(), 300); });
downloadBtn.addEventListener('click', downloadCSV);

$('btn-watch').addEventListener('click', () => setMode(false));
$('btn-play').addEventListener('click',  () => setMode(true));
$('play-stay-btn').addEventListener('click',   () => onPlayChoice('stay'));
$('play-switch-btn').addEventListener('click', () => onPlayChoice('switch'));

mathBtn.addEventListener('click', () => {
  const open = mathBtn.getAttribute('aria-expanded') === 'true';
  mathBtn.setAttribute('aria-expanded', !open);
  mathBody.setAttribute('aria-hidden', open);
  mathBody.classList.toggle('open');
});

// ============================================================
// Init
// ============================================================
buildDoors(numDoors);
updateTheoDisplay();
renderMath(numDoors);
// Theme init happens here so chart already exists when updateChartTheme is called
const saved = localStorage.getItem('theme');
setTheme(saved !== 'light');
setTimeout(() => runDemo(), 600);
