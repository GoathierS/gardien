// Copie l'interface commune (web/static) dans www, le dossier que Capacitor embarque.
// web/static reste la seule source : ne pas modifier www à la main.
import { cpSync, rmSync } from "node:fs";
import { fileURLToPath } from "node:url";

const source = fileURLToPath(new URL("../web/static", import.meta.url));
const cible = fileURLToPath(new URL("./www", import.meta.url));

rmSync(cible, { recursive: true, force: true });
cpSync(source, cible, { recursive: true });
console.log(`Interface copiée : ${source} -> ${cible}`);
