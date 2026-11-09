// Interface commune de gardien : appli de bureau (Wails), page servie par un agent
// (PC, NAS) et appli téléphone (Capacitor).
//
// L'API est relative ("/api/v1/...") quand l'interface est servie par un agent.
// Dans l'appli téléphone, il faudra indiquer l'adresse d'un agent du réseau.

"use strict";

const SurTelephone = Boolean(window.Capacitor?.isNativePlatform?.());
const ClefAgent = "gardien.agent";

function adresseAgent() {
  return SurTelephone ? (localStorage.getItem(ClefAgent) || "") : "";
}

async function api(chemin) {
  const reponse = await fetch(adresseAgent() + "/api/v1" + chemin);
  if (!reponse.ok) throw new Error(`HTTP ${reponse.status}`);
  return reponse.json();
}

const libelles = { ok: "OK", attention: "Attention", critique: "Critique", inconnu: "Inconnu" };
const gravite = { ok: 0, inconnu: 1, attention: 2, critique: 3 };
const enGo = (o) => (o / 1024 ** 3).toLocaleString("fr-FR", { maximumFractionDigits: 1 }) + " Go";

function element(balise, classe, texte) {
  const e = document.createElement(balise);
  if (classe) e.className = classe;
  if (texte !== undefined) e.textContent = texte;
  return e;
}

function pastille(niveau) {
  return element("span", "pastille " + niveau, libelles[niveau] ?? niveau);
}

// Rendu spécifique par capteur. Un capteur sans rendu dédié n'affiche que son résumé.
const rendus = {
  "espace-disque"(details) {
    const table = element("table");
    const entete = table.createTHead().insertRow();
    for (const [texte, classe] of [["Volume"], ["Occupation"], ["Libre", "nombre"], ["Total", "nombre"]]) {
      entete.append(element("th", classe, texte));
    }
    const corps = table.createTBody();
    for (const v of details ?? []) {
      const ligne = corps.insertRow();
      ligne.insertCell().textContent = v.point;

      const occupation = ligne.insertCell();
      const barre = element("div", "barre " + v.niveau);
      const remplissage = element("span");
      remplissage.style.width = Math.min(100, v.utilise_pourcent) + "%";
      barre.append(remplissage);
      occupation.append(barre, document.createTextNode(Math.round(v.utilise_pourcent) + " %"));

      ligne.insertCell().textContent = enGo(v.libre_octets);
      ligne.insertCell().textContent = enGo(v.total_octets);
      ligne.cells[2].className = ligne.cells[3].className = "nombre";
    }
    return table;
  },
};

function carte(rapport) {
  const c = element("section", "carte");
  const entete = element("div", "entete-carte");
  entete.append(element("h2", "", rapport.nom), pastille(rapport.niveau));
  c.append(entete, element("p", "resume", rapport.resume));
  const rendu = rendus[rapport.capteur];
  if (rendu && rapport.details) c.append(rendu(rapport.details));
  return c;
}

async function charger() {
  const zone = document.getElementById("rapports");
  const global = document.getElementById("etat-global");
  try {
    const [version, rapports] = await Promise.all([api("/version"), api("/sante")]);
    document.getElementById("version").textContent = version.version;

    const pire = rapports.reduce((n, r) => (gravite[r.niveau] > gravite[n] ? r.niveau : n), "ok");
    global.className = "pastille " + pire;
    global.textContent = libelles[pire];

    zone.replaceChildren(...rapports.map(carte));
  } catch {
    global.className = "pastille inconnu";
    global.textContent = "Hors ligne";
    zone.replaceChildren(element("p", "erreur",
      SurTelephone
        ? "Aucun agent gardien joignable. La connexion du téléphone à un PC ou un NAS arrive dans la v0.4."
        : "Impossible de joindre l'agent gardien."));
  }
}

if (SurTelephone && !adresseAgent()) {
  const bandeau = document.getElementById("bandeau");
  bandeau.textContent = "Version téléphone en préparation : elle affichera l'état de tes PC et de ton NAS.";
  bandeau.hidden = false;
}

charger();
setInterval(charger, 30000);
