import { View, StyleSheet } from 'react-native';
import Svg, { Ellipse, Path } from 'react-native-svg';
import { roastColor } from '@/src/api/beans';

interface Props {
  roastId?: string | null;
  size: number;
}

export default function RoastBubble({ roastId, size }: Props) {
  const bgColor = roastColor(roastId);
  const borderRadius = Math.round(size * 0.38);
  const iconW = Math.round(size * 0.46);
  const iconH = Math.round(iconW * 28 / 24);

  return (
    <View style={[styles.bubble, { width: size, height: size, borderRadius, backgroundColor: bgColor }]}>
      <Svg width={iconW} height={iconH} viewBox="0 0 24 28">
        <Ellipse cx={12} cy={14} rx={9} ry={12} fill="#FDF8F2" />
        <Path
          d="M12 3 C14.5 8.5 14.5 19.5 12 25"
          stroke={bgColor}
          strokeWidth={1.8}
          strokeLinecap="round"
          fill="none"
        />
      </Svg>
    </View>
  );
}

const styles = StyleSheet.create({
  bubble: { alignItems: 'center', justifyContent: 'center' },
});
