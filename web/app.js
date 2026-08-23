"use strict";

const compInput = document.getElementById("comp");
const tempInput = document.getElementById("temp");
const resultBox = document.getElementById("result");
const errorBox = document.getElementById("error");
const exampleNote = document.getElementById("example-note");
const scanMeta = document.getElementById("scan-meta");
const chart = document.getElementById("chart");

function showError(msg) {
  errorBox.textContent = msg;
}

function clearError() {
  errorBox.textContent = "";
}

function showResult(text) {
  resultBox.textContent = text;
}

async function postJSON(path, body) {
  const resp = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await resp.json();
  if (!resp.ok) {
    throw new Error(data.error || ("HTTP " + resp.status));
  }
  return data;
}

function formatPoint(data) {
  const lines = [];
  lines.push("相区: " + data.region + "  —  " + data.region_label);
  for (const ph of data.phases) {
    lines.push(
      "  " + ph.name + ": 成分 c = " + ph.composition.toFixed(4) +
        ", 质量分数 w = " + ph.mass_fraction.toFixed(4)
    );
  }
  lines.push("液相分数 fL = " + data.liquid_fraction.toFixed(4));
  if (typeof data.eutectic_fraction === "number") {
    lines.push("刚低于共晶温度时凝固的共晶组织分数 = " + data.eutectic_fraction.toFixed(4));
  }
  lines.push("两相质量分数之和 = " + data.total_fraction.toFixed(6) + "（应为 1）");
  lines.push("说明: " + data.note);
  return lines.join("\n");
}

async function computePoint() {
  clearError();
  const c = Number(compInput.value);
  const t = Number(tempInput.value);
  try {
    const data = await postJSON("/api/point", { c: c, t: t });
    showResult(formatPoint(data));
  } catch (err) {
    showResult("");
    showError("API 错误: " + err.message);
  }
}

async function scanCurve() {
  clearError();
  const c = Number(compInput.value);
  const tMin = Number(tempInput.value) - 180;
  const tMax = Number(tempInput.value) + 120;
  try {
    const data = await postJSON("/api/scan", {
      c: c,
      tmin: tMin < 0 ? 0 : tMin,
      tmax: tMax,
      n: 121,
    });
    drawScan(data);
    scanMeta.textContent =
      "扫描成分 c = " + data.composition.toFixed(4) +
      "，共晶温度 TE = " + data.teutectic_k.toFixed(2) + " K，点数 = " + data.points.length;
  } catch (err) {
    showError("扫描 API 错误: " + err.message);
  }
}

function drawScan(data) {
  const W = 680, H = 320, padL = 46, padR = 14, padT = 14, padB = 34;
  const pts = data.points;
  let minT = Infinity, maxT = -Infinity;
  for (const p of pts) {
    if (p.temperature < minT) minT = p.temperature;
    if (p.temperature > maxT) maxT = p.temperature;
  }
  if (maxT <= minT) maxT = minT + 1;
  const x = (t) => padL + ((t - minT) / (maxT - minT)) * (W - padL - padR);
  const y = (f) => padT + (1 - f) * (H - padT - padB);

  let svg = "";
  // grid + axes
  svg += `<line x1="${padL}" y1="${padT}" x2="${padL}" y2="${H - padB}" stroke="#c8d0d8"/>`;
  svg += `<line x1="${padL}" y1="${H - padB}" x2="${W - padR}" y2="${H - padB}" stroke="#c8d0d8"/>`;
  for (let i = 0; i <= 4; i++) {
    const f = i / 4;
    svg += `<line x1="${padL}" y1="${y(f)}" x2="${W - padR}" y2="${y(f)}" stroke="#eef1f4"/>`;
    svg += `<text x="${padL - 6}" y="${y(f) + 4}" text-anchor="end" font-size="11" fill="#5c6b7a">${f.toFixed(2)}</text>`;
  }
  svg += `<text x="${padL}" y="${padT - 6}" font-size="11" fill="#5c6b7a">fL</text>`;
  svg += `<text x="${padL}" y="${H - 8}" font-size="11" fill="#5c6b7a">${minT.toFixed(0)} K</text>`;
  svg += `<text x="${W - padR}" y="${H - 8}" text-anchor="end" font-size="11" fill="#5c6b7a">${maxT.toFixed(0)} K</text>`;
  svg += `<text x="${padL + 10}" y="${padT + 10}" font-size="11" fill="#5c6b7a">fL(T) 曲线，点列来自 /api/scan</text>`;

  // polyline from the server points
  let d = "";
  for (let i = 0; i < pts.length; i++) {
    const px = x(pts[i].temperature);
    const py = y(pts[i].liquid_fraction);
    d += (i === 0 ? "M" : "L") + px.toFixed(2) + " " + py.toFixed(2);
  }
  svg += `<path d="${d}" fill="none" stroke="#0f6fde" stroke-width="2" stroke-linejoin="round"/>`;

  // eutectic temperature reference line
  if (data.teutectic_k > minT && data.teutectic_k < maxT) {
    const ex = x(data.teutectic_k);
    svg += `<line x1="${ex}" y1="${padT}" x2="${ex}" y2="${H - padB}" stroke="#b3261e" stroke-dasharray="4 3"/>`;
    svg += `<text x="${ex + 3}" y="${padT + 10}" font-size="10" fill="#b3261e">TE</text>`;
  }

  chart.innerHTML = svg;
}

document.getElementById("load-example").addEventListener("click", async () => {
  clearError();
  try {
    const resp = await fetch("/api/examples");
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || ("HTTP " + resp.status));
    const ex = data["alloy-60"];
    compInput.value = ex.c;
    tempInput.value = ex.t;
    exampleNote.textContent = ex.label + "  |  " + ex.note;
  } catch (err) {
    showError("加载示例错误: " + err.message);
  }
});

document.getElementById("compute").addEventListener("click", computePoint);
document.getElementById("scan").addEventListener("click", scanCurve);

computePoint();
