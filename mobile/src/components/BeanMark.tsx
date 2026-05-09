import Svg, { Defs, RadialGradient, Stop, Path } from 'react-native-svg';

export function BeanMark({ size = 48 }: { size?: number }) {
  const h = size * 1.22;
  return (
    <Svg width={size} height={h} viewBox="0 0 36 44">
      <Defs>
        <RadialGradient id="bmbg" cx="38%" cy="30%" r="65%" gradientUnits="userSpaceOnUse" fx="38%" fy="30%">
          <Stop offset="0%" stopColor="#3A1E10" />
          <Stop offset="100%" stopColor="#0F0603" />
        </RadialGradient>
      </Defs>
      <Path d="M9 9 C7.5 6 10 3.5 8.5 1"   stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      <Path d="M18 7 C16.5 4 19 1.5 17.5 -1" stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      <Path d="M27 9 C25.5 6 28 3.5 26.5 1" stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      <Path
        d="M13 12 C6 14 2 21 2 28 C2 36 7 44 14 44 C17 44.5 19 44.5 22 44 C29 43 34 36 34 28 C34 20 29 13 22 12 C19 11 16 11 13 12 Z"
        fill="url(#bmbg)"
      />
      <Path d="M21 13 C15 23 23 31 15 40 C13 43 16 44 18 44" stroke="#FAF3E8" strokeWidth="1.5" strokeLinecap="round" opacity="0.5" />
    </Svg>
  );
}
