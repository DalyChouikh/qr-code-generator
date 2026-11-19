export const THEME_STORAGE_KEY = "qrgen-theme";

// Runs inline in <head> before first paint so a saved light theme doesn't flash dark.
export const themeInitScript = `try{if(localStorage.getItem("${THEME_STORAGE_KEY}")==="light")document.documentElement.dataset.theme="light"}catch(e){}`;
