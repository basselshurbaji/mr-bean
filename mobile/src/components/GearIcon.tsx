import Svg, { Circle, Line, Path, Rect } from 'react-native-svg';

interface Props {
  typeId: string;
  size?: number;
  color?: string;
}

const SW = 1.8;
const SC = 'round';
const SJ = 'round';

export default function GearIcon({ typeId, size = 24, color = '#1C0F07' }: Props) {
  const props = { stroke: color, strokeWidth: SW, strokeLinecap: SC as 'round', strokeLinejoin: SJ as 'round', fill: 'none' };

  switch (typeId) {
    case 'machine':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Rect x="2" y="8" width="20" height="12" rx="2.5" {...props} />
          <Rect x="7" y="18" width="10" height="3" rx="1" {...props} />
          <Line x1="2" y1="12" x2="22" y2="12" {...props} />
          <Circle cx="6.5" cy="10" r="1.5" {...props} />
          <Circle cx="17.5" cy="10" r="1" {...props} />
          <Path d="M22 10 Q25 10 25 14" {...props} />
          <Path d="M21 14 L23 17" {...props} />
        </Svg>
      );
    case 'grinder':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Path d="M8.5 10 L10.5 3.5 L13.5 3.5 L15.5 10 Z" {...props} />
          <Rect x="7" y="10" width="10" height="9" rx="2" {...props} />
          <Line x1="7" y1="13.5" x2="17" y2="13.5" {...props} />
          <Rect x="10" y="19" width="4" height="2.5" rx="1" {...props} />
        </Svg>
      );
    case 'scale':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Rect x="2" y="17" width="20" height="4" rx="2" {...props} />
          <Rect x="4" y="13.5" width="16" height="4" rx="1.5" {...props} />
          <Rect x="6" y="7" width="12" height="6" rx="1.5" {...props} />
          <Line x1="9" y1="10" x2="15" y2="10" {...props} />
        </Svg>
      );
    case 'portafilter':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Circle cx="12" cy="10" r="8" {...props} />
          <Circle cx="12" cy="10" r="6" {...props} strokeWidth={1} opacity={0.4} />
          <Path d="M7 7 A6 6 0 0 1 12 4" {...props} strokeWidth={1.5} opacity={0.5} />
          <Circle cx="12" cy="10" r="1" fill={color} stroke="none" />
          <Circle cx="9" cy="8" r="0.8" fill={color} stroke="none" opacity={0.9} />
          <Circle cx="15" cy="8" r="0.8" fill={color} stroke="none" opacity={0.9} />
          <Circle cx="9" cy="12" r="0.8" fill={color} stroke="none" opacity={0.9} />
          <Circle cx="15" cy="12" r="0.8" fill={color} stroke="none" opacity={0.9} />
          <Circle cx="12" cy="6.2" r="0.7" fill={color} stroke="none" opacity={0.7} />
          <Circle cx="12" cy="13.8" r="0.7" fill={color} stroke="none" opacity={0.7} />
          <Line x1="10" y1="18" x2="10" y2="23" {...props} strokeWidth={2} />
          <Line x1="14" y1="18" x2="14" y2="23" {...props} strokeWidth={2} />
        </Svg>
      );
    case 'tamper':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Path d="M9.5 2 C9.5 1 14.5 1 14.5 2 L14.5 7 C14.5 8 13 9 12 9 C11 9 9.5 8 9.5 7 Z" {...props} />
          <Rect x="10.5" y="9" width="3" height="5" {...props} />
          <Rect x="6" y="14" width="12" height="5" rx="1.5" {...props} />
        </Svg>
      );
    case 'distributor':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Circle cx="12" cy="14" r="7.5" {...props} />
          <Line x1="12" y1="6.5" x2="12" y2="14" {...props} />
          <Line x1="18.5" y1="10.5" x2="12" y2="14" {...props} />
          <Line x1="18.5" y1="17.5" x2="12" y2="14" {...props} />
          <Line x1="12" y1="21.5" x2="12" y2="14" {...props} />
          <Line x1="5.5" y1="17.5" x2="12" y2="14" {...props} />
          <Line x1="5.5" y1="10.5" x2="12" y2="14" {...props} />
          <Circle cx="12" cy="14" r="2" fill={color} stroke="none" />
          <Line x1="12" y1="2" x2="12" y2="6.5" {...props} strokeWidth={2.5} />
        </Svg>
      );
    case 'wdt':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Path d="M10 2 C10 1 14 1 14 2 L14 11 L10 11 Z" {...props} />
          <Line x1="10" y1="11" x2="5" y2="22" {...props} />
          <Line x1="11" y1="11" x2="9" y2="22" {...props} />
          <Line x1="12" y1="11" x2="12" y2="22" {...props} />
          <Line x1="13" y1="11" x2="15" y2="22" {...props} />
          <Line x1="14" y1="11" x2="19" y2="22" {...props} />
        </Svg>
      );
    case 'basket':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Line x1="3" y1="7" x2="21" y2="7" {...props} />
          <Path d="M4 7 L6.5 19 C6.5 20.5 17.5 20.5 17.5 19 L20 7" {...props} />
          <Circle cx="9" cy="13" r="0.9" fill={color} stroke="none" />
          <Circle cx="12" cy="13" r="0.9" fill={color} stroke="none" />
          <Circle cx="15" cy="13" r="0.9" fill={color} stroke="none" />
          <Circle cx="10.5" cy="16.5" r="0.9" fill={color} stroke="none" />
          <Circle cx="13.5" cy="16.5" r="0.9" fill={color} stroke="none" />
        </Svg>
      );
    case 'puckscreen':
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Circle cx="12" cy="12" r="9" {...props} />
          <Circle cx="12" cy="12" r="4.5" {...props} />
          <Circle cx="12" cy="12" r="1.2" fill={color} stroke="none" />
          <Circle cx="12" cy="7" r="0.8" fill={color} stroke="none" />
          <Circle cx="16.5" cy="9" r="0.8" fill={color} stroke="none" />
          <Circle cx="16.5" cy="15" r="0.8" fill={color} stroke="none" />
          <Circle cx="12" cy="17" r="0.8" fill={color} stroke="none" />
          <Circle cx="7.5" cy="15" r="0.8" fill={color} stroke="none" />
          <Circle cx="7.5" cy="9" r="0.8" fill={color} stroke="none" />
        </Svg>
      );
    default:
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24">
          <Circle cx="12" cy="12" r="9" {...props} />
          <Line x1="12" y1="8" x2="12" y2="13" {...props} />
          <Circle cx="12" cy="16" r="0.8" fill={color} stroke="none" />
        </Svg>
      );
  }
}
