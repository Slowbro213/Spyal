import { readdirSync, readFileSync, writeFileSync, statSync } from 'fs';
import { join, extname } from 'path';
import { brotliCompressSync } from 'node:zlib';

// Directory with built static files
const outDir = './public/out';

// File types to compress
const compressExts = ['.js', '.css', '.html'];

// Function to recursively compress files
async function compressFiles(dir: string) {
  try {
    for (const file of readdirSync(dir)) {
      const filePath = join(dir, file);
      const stat = statSync(filePath);

      if (stat.isDirectory()) {
        await compressFiles(filePath); // Recursively process subdirectories
        continue;
      }

      const ext = extname(file);
      if (!compressExts.includes(ext) || file.endsWith('.br')) continue;

      try {
        const data = readFileSync(filePath);
        const compressed = brotliCompressSync(data);

        const outPath = filePath + '.br';
        writeFileSync(outPath, compressed);

        console.log(`✅ Compressed ${filePath} → ${outPath}`);
      } catch (err) {
        console.error(`❌ Failed to compress ${filePath}: ${err}`);
      }
    }
  } catch (err) {
    console.error(`❌ Error reading directory ${dir}: ${err}`);
  }
}

console.log(`Starting Brotli compression for files in ${outDir}`);
compressFiles(outDir).then(() => {
  console.log('Compression complete');
});
