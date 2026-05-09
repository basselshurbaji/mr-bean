import { Animated, Pressable, StyleSheet, View } from 'react-native';
import { KeyboardAvoidingView } from 'react-native-keyboard-controller';
import { palette, radii } from '@/src/theme';

interface Props {
  backdropOpacity: Animated.Value;
  translateY: Animated.Value;
  onClose: () => void;
  children: React.ReactNode;
}

export function SheetLayout({ backdropOpacity, translateY, onClose, children }: Props) {
  return (
    <View style={StyleSheet.absoluteFillObject}>
      {/*
       * Backdrop is a sibling of KAV — not a child.
       * This prevents it from being clipped when KAV shrinks its content
       * area on keyboard show, which causes the visible flash/glitch.
       */}
      <Animated.View style={[StyleSheet.absoluteFillObject, styles.backdrop, { opacity: backdropOpacity }]}>
        <Pressable style={StyleSheet.absoluteFill} onPress={onClose} />
      </Animated.View>

      {/*
       * KAV pushes the sheet above the keyboard (behavior="padding").
       * KeyboardAwareScrollView inside the sheet then auto-scrolls to the
       * focused field within the lifted sheet.
       */}
      <KeyboardAvoidingView style={styles.kav} behavior="padding" pointerEvents="box-none">
        <Animated.View style={[styles.sheet, { transform: [{ translateY }] }]}>
          <View style={styles.handle} />
          {children}
        </Animated.View>
      </KeyboardAvoidingView>
    </View>
  );
}

const styles = StyleSheet.create({
  backdrop: {
    backgroundColor: 'rgba(28,15,7,0.42)',
  },
  kav: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  sheet: {
    maxHeight: '92%',
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.16,
    shadowRadius: 32,
    elevation: 12,
  },
  handle: {
    width: 36,
    height: 4,
    borderRadius: radii.full,
    backgroundColor: palette.cream500,
    alignSelf: 'center',
    marginTop: 14,
    marginBottom: 4,
  },
});
