export const fontFamilies = {
  display: 'PlayfairDisplay_700Bold',
  displayItalic: 'PlayfairDisplay_700Bold_Italic',
  body: 'DMSans_400Regular',
  bodyMedium: 'DMSans_500Medium',
  bodySemiBold: 'DMSans_700Bold',
  mono: 'JetBrainsMono_400Regular',
  monoMedium: 'JetBrainsMono_500Medium',
} as const;

export const fontSizes = {
  xs: 11,
  sm: 13,
  base: 15,
  md: 17,
  lg: 20,
  xl: 24,
  '2xl': 30,
  '3xl': 38,
  '4xl': 48,
  '5xl': 64,
} as const;

export const lineHeights = {
  tight: 1.15,
  snug: 1.3,
  normal: 1.5,
  loose: 1.7,
} as const;

export const letterSpacings = {
  tight: -0.02,
  normal: 0,
  wide: 0.04,
  wider: 0.08,
} as const;

export const textStyles = {
  h1: {
    fontFamily: fontFamilies.display,
    fontSize: fontSizes['4xl'],
    lineHeight: fontSizes['4xl'] * lineHeights.tight,
    letterSpacing: fontSizes['4xl'] * letterSpacings.tight,
  },
  h2: {
    fontFamily: fontFamilies.display,
    fontSize: fontSizes['3xl'],
    lineHeight: fontSizes['3xl'] * lineHeights.tight,
    letterSpacing: fontSizes['3xl'] * letterSpacings.tight,
  },
  h3: {
    fontFamily: fontFamilies.bodySemiBold,
    fontSize: fontSizes.xl,
    lineHeight: fontSizes.xl * lineHeights.snug,
  },
  h4: {
    fontFamily: fontFamilies.bodySemiBold,
    fontSize: fontSizes.lg,
    lineHeight: fontSizes.lg * lineHeights.snug,
  },
  body: {
    fontFamily: fontFamilies.body,
    fontSize: fontSizes.base,
    lineHeight: fontSizes.base * lineHeights.normal,
  },
  bodySm: {
    fontFamily: fontFamilies.body,
    fontSize: fontSizes.sm,
    lineHeight: fontSizes.sm * lineHeights.normal,
  },
  caption: {
    fontFamily: fontFamilies.body,
    fontSize: fontSizes.xs,
    lineHeight: fontSizes.xs * lineHeights.snug,
    letterSpacing: fontSizes.xs * letterSpacings.wide,
    textTransform: 'uppercase' as const,
  },
  label: {
    fontFamily: fontFamilies.bodyMedium,
    fontSize: fontSizes.sm,
  },
  mono: {
    fontFamily: fontFamilies.mono,
    fontSize: fontSizes.base,
  },
  monoLg: {
    fontFamily: fontFamilies.monoMedium,
    fontSize: fontSizes['2xl'],
  },
  monoXl: {
    fontFamily: fontFamilies.monoMedium,
    fontSize: fontSizes['4xl'],
  },
} as const;
