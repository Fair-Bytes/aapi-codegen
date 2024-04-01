import { defineConfig } from "tsup";
 
export default defineConfig({
  entry: ["src/aapi-codegen.ts"],
  publicDir: false,
  clean: true,
  minify: true,
  format: ["cjs"],
});