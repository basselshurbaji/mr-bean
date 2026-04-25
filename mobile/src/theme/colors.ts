export const palette = {
  espresso900: '#0D0603',
  espresso800: '#1C0F07',
  espresso700: '#2D1810',
  espresso600: '#4A2518',
  espresso500: '#6B3A2A',
  espresso400: '#8C5340',
  espresso300: '#B07060',

  caramel600: '#9B5E1A',
  caramel500: '#C4782A',
  caramel400: '#D4872A',
  caramel300: '#E4A855',
  caramel200: '#F0C882',
  caramel100: '#F8E4B8',

  cream600: '#C4A882',
  cream500: '#D8C4A0',
  cream400: '#E8D8B8',
  cream300: '#F0E4CC',
  cream200: '#FAF3E8',
  cream100: '#FDF8F2',

  matcha700: '#2D5235',
  matcha600: '#3A6644',
  matcha500: '#4A7C59',
  matcha400: '#6A9B74',
  matcha300: '#92BC9C',
  matcha200: '#C0D9C4',
  matcha100: '#E8F2EA',

  error600: '#9B2318',
  error500: '#C0392B',
  error100: '#FDECEA',
  warning500: '#D68910',
  warning100: '#FEF5DC',
  success500: '#4A7C59',
  success100: '#E8F2EA',
} as const;

export const colors = {
  bgApp: palette.cream100,
  bgCard: palette.cream200,
  bgSubtle: palette.cream300,
  bgInverse: palette.espresso800,
  bgAccent: palette.caramel400,

  fgPrimary: palette.espresso800,
  fgSecondary: palette.espresso500,
  fgTertiary: palette.espresso300,
  fgDisabled: palette.cream500,
  fgInverse: palette.cream100,
  fgAccent: palette.caramel500,
  fgLink: palette.matcha500,

  borderSubtle: palette.cream400,
  borderDefault: palette.cream500,
  borderStrong: palette.espresso400,
  borderFocus: palette.caramel400,

  interactivePrimary: palette.espresso800,
  interactivePrimaryHover: palette.espresso700,
  interactiveSecondary: palette.caramel400,
  interactiveSecondaryHover: palette.caramel500,
  interactiveAccent: palette.matcha500,
  interactiveAccentHover: palette.matcha600,
  interactiveDestructive: palette.error500,

  // Dark mode
  dark: {
    bgApp: '#110906',
    bgCard: '#1E0F08',
    bgSubtle: '#2A160C',
    fgPrimary: palette.cream100,
    fgSecondary: palette.cream500,
    fgTertiary: '#8C6050',
    border: '#3A200E',
  },
} as const;
