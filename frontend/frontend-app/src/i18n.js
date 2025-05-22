// Simple in-app i18n
const messages = {
  en: {
    title: 'QRazem',
    textOrUrl: 'Text or URL',
    enterTextUrl: 'Enter text or URL',
    qrColorLabel: 'QR Color',
    bgColorLabel: 'Background Color',
    chooseColor: 'Choose color',
    addLogo: 'Add a logo',
    svgOnly: 'SVG only',
    uploadSvg: 'Upload logo',
    orDragDrop: 'or drag and drop here',
    searchLogos: 'Search logos...',
    generateQr: 'Generate QR',
    preview: 'Preview',
    download: 'Download',
    imageOnly: 'Allowed formats: SVG, PNG, JPEG',
    failed: 'Failed to generate QR code.'
  },
  pl: {
    title: 'QRazem',
    textOrUrl: 'Tekst lub adres URL',
    enterTextUrl: 'Wpisz tekst lub adres URL',
    qrColorLabel: 'Kolor QR',
    bgColorLabel: 'Kolor tła',
    chooseColor: 'Wybierz kolor',
    addLogo: 'Dodaj logo',
    svgOnly: 'Tylko SVG',
    uploadSvg: 'Prześlij logo',
    orDragDrop: 'lub przeciągnij i upuść tutaj',
    searchLogos: 'Wyszukaj logo...',
    generateQr: 'Generuj QR',
    preview: 'Podgląd',
    download: 'Pobierz',
    imageOnly: 'Dozwolone formaty: SVG, PNG, JPEG',
    failed: 'Nie udało się wygenerować kodu QR.'
  }
};

// Detect browser language
const locale =
  typeof navigator !== 'undefined' && navigator.language.startsWith('pl')
    ? 'pl'
    : 'en';

/**
 * Translate a key using in-app messages
 * @param {string} key
 * @returns {string}
 */
export function t(key) {
  return messages[locale][key] || messages['pl'][key] || key;
}