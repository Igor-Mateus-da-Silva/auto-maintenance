import fs from 'node:fs';
import path from 'node:path';

const targetPath = path.resolve('src/environments/environment.prod.ts');
const apiUrl = (process.env.API_URL || 'http://localhost:8080/api').trim();

const content = `export const environment = {
  production: true,
  apiUrl: ${JSON.stringify(apiUrl)}
};
`;

fs.mkdirSync(path.dirname(targetPath), { recursive: true });
fs.writeFileSync(targetPath, content, 'utf8');
console.log(`[set-env] environment.prod.ts atualizado com apiUrl=${apiUrl}`);
