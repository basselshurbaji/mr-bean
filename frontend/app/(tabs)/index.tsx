import { View, Text, StyleSheet } from 'react-native';
import { colors, textStyles, spacing } from '@/src/theme';

export default function HomeScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.heading}>Home</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.bgApp,
    padding: spacing[4],
  },
  heading: {
    ...textStyles.h3,
    color: colors.fgPrimary,
  },
});
