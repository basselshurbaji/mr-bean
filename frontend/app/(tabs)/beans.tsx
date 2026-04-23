import { View, Text, StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, textStyles, spacing } from '@/src/theme';

export default function BeansScreen() {
  return (
    <SafeAreaView style={styles.container} edges={['top']}>
      <Text style={styles.heading}>Your beans</Text>
    </SafeAreaView>
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
