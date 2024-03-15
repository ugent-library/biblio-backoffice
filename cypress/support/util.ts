export function getRandomText() {
  return crypto.randomUUID().replace(/-/g, "").toUpperCase();
}
