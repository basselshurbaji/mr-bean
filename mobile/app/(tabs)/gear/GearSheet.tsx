import {
  ActivityIndicator,
  Animated,
  KeyboardAvoidingView,
  Modal,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { useEffect, useRef, useState } from 'react';
import { colors, palette, radii, spacing } from '@/src/theme';
import { gearApi, GearItem } from '@/src/api/gear';
import GearIcon from '@/src/components/GearIcon';

const TYPES = [
  { id: 'machine',     label: 'Espresso machine'  },
  { id: 'grinder',     label: 'Grinder'           },
  { id: 'scale',       label: 'Scale'             },
  { id: 'portafilter', label: 'Portafilter'       },
  { id: 'tamper',      label: 'Tamper'            },
  { id: 'distributor', label: 'Distribution tool' },
  { id: 'wdt',         label: 'WDT tool'          },
  { id: 'basket',      label: 'Basket'            },
  { id: 'puckscreen',  label: 'Puck screen'       },
  { id: 'other',       label: 'Other'             },
];

interface Props {
  editItem?: GearItem;
  onClose: () => void;
  onSaved: (item: GearItem) => void;
}

export default function GearSheet({ editItem, onClose, onSaved }: Props) {
  const translateY = useRef(new Animated.Value(600)).current;
  const backdropOpacity = useRef(new Animated.Value(0)).current;

  const isEdit = !!editItem;
  const [step, setStep] = useState<'type' | 'form'>(isEdit ? 'form' : 'type');
  const [selectedType, setSelectedType] = useState(editItem?.type_id ?? '');
  const [form, setForm] = useState({
    name: editItem?.name ?? '',
    brand: editItem?.brand ?? '',
    model: editItem?.model ?? '',
    year: editItem?.year ?? '',
    notes: editItem?.notes ?? '',
  });
  const [focused, setFocused] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Animated.parallel([
      Animated.timing(translateY, { toValue: 0, duration: 300, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 1, duration: 300, useNativeDriver: true }),
    ]).start();
  }, []);

  function dismiss() {
    Animated.parallel([
      Animated.timing(translateY, { toValue: 600, duration: 260, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 0, duration: 260, useNativeDriver: true }),
    ]).start(() => onClose());
  }

  async function save() {
    if (saving || !form.name.trim()) return;
    setSaving(true);
    setError(null);
    try {
      const body = {
        type_id: selectedType,
        name: form.name.trim(),
        brand: form.brand.trim() || undefined,
        model: form.model.trim() || undefined,
        year: form.year.trim() || undefined,
        notes: form.notes.trim() || undefined,
      };
      const result = isEdit
        ? await gearApi.updateGear(editItem!.id, body)
        : await gearApi.createGear(body);
      onSaved(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went wrong.');
      setSaving(false);
    }
  }

  const typeLabel = TYPES.find(t => t.id === selectedType)?.label ?? '';
  const canSave = form.name.trim().length > 0;

  return (
    <Modal transparent animationType="none" onRequestClose={dismiss}>
      <View style={styles.overlay}>
        <Animated.View
          style={[styles.backdrop, { opacity: backdropOpacity }]}
        >
          <Pressable style={StyleSheet.absoluteFill} onPress={dismiss} />
        </Animated.View>

        <Animated.View style={[styles.sheet, { transform: [{ translateY }] }]}>
          {/* Drag handle */}
          <View style={styles.handle} />

          <KeyboardAvoidingView
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
          >
            {step === 'type' ? (
              <ScrollView
                contentContainerStyle={styles.typeContent}
                showsVerticalScrollIndicator={false}
              >
                <Text style={styles.sheetTitle}>What type of gear?</Text>
                <View style={styles.typeGrid}>
                  {TYPES.map(t => (
                    <Pressable
                      key={t.id}
                      style={({ pressed }) => [
                        styles.typeBtn,
                        pressed && styles.typeBtnPressed,
                      ]}
                      onPress={() => {
                        setSelectedType(t.id);
                        setStep('form');
                      }}
                    >
                      <GearIcon typeId={t.id} size={30} color={palette.espresso800} />
                      <Text style={styles.typeBtnLabel}>{t.label}</Text>
                    </Pressable>
                  ))}
                </View>
              </ScrollView>
            ) : (
              <ScrollView
                contentContainerStyle={styles.formContent}
                showsVerticalScrollIndicator={false}
                keyboardShouldPersistTaps="handled"
              >
                {!isEdit && (
                  <Pressable onPress={() => setStep('type')} style={styles.backLink}>
                    <Text style={styles.backLinkText}>← Change type</Text>
                  </Pressable>
                )}
                <Text style={styles.sheetTitle}>
                  {isEdit ? 'Edit gear' : 'Add gear'}
                </Text>
                <Text style={styles.sheetSub}>{typeLabel}</Text>

                {/* Name */}
                <View style={styles.fieldWrap}>
                  <Text style={styles.fieldLabel}>
                    Name <Text style={styles.required}>*</Text>
                  </Text>
                  <TextInput
                    style={[styles.input, focused === 'name' && styles.inputFocused]}
                    placeholder="e.g. Niche Zero"
                    placeholderTextColor={colors.fgDisabled}
                    value={form.name}
                    onChangeText={v => setForm(f => ({ ...f, name: v }))}
                    onFocus={() => setFocused('name')}
                    onBlur={() => setFocused(null)}
                    returnKeyType="next"
                    textContentType="none"
                    autoComplete="off"
                  />
                </View>

                {/* Brand + Model row */}
                <View style={styles.row}>
                  <View style={styles.rowCol}>
                    <Text style={styles.fieldLabel}>Brand</Text>
                    <TextInput
                      style={[styles.input, focused === 'brand' && styles.inputFocused]}
                      placeholder="Niche"
                      placeholderTextColor={colors.fgDisabled}
                      value={form.brand}
                      onChangeText={v => setForm(f => ({ ...f, brand: v }))}
                      onFocus={() => setFocused('brand')}
                      onBlur={() => setFocused(null)}
                      returnKeyType="next"
                      textContentType="none"
                      autoComplete="off"
                    />
                  </View>
                  <View style={styles.rowCol}>
                    <Text style={styles.fieldLabel}>Model</Text>
                    <TextInput
                      style={[styles.input, focused === 'model' && styles.inputFocused]}
                      placeholder="Zero"
                      placeholderTextColor={colors.fgDisabled}
                      value={form.model}
                      onChangeText={v => setForm(f => ({ ...f, model: v }))}
                      onFocus={() => setFocused('model')}
                      onBlur={() => setFocused(null)}
                      returnKeyType="next"
                      textContentType="none"
                      autoComplete="off"
                    />
                  </View>
                </View>

                {/* Year */}
                <View style={styles.fieldWrap}>
                  <Text style={styles.fieldLabel}>Year purchased</Text>
                  <TextInput
                    style={[styles.input, focused === 'year' && styles.inputFocused]}
                    placeholder="2022"
                    placeholderTextColor={colors.fgDisabled}
                    value={form.year}
                    onChangeText={v => setForm(f => ({ ...f, year: v.replace(/\D/g, '').slice(0, 4) }))}
                    onFocus={() => setFocused('year')}
                    onBlur={() => setFocused(null)}
                    keyboardType="number-pad"
                    maxLength={4}
                    returnKeyType="next"
                    textContentType="none"
                    autoComplete="off"
                  />
                </View>

                {/* Notes */}
                <View style={styles.fieldWrap}>
                  <Text style={styles.fieldLabel}>Notes</Text>
                  <TextInput
                    style={[styles.input, styles.textarea, focused === 'notes' && styles.inputFocused]}
                    placeholder="Any details worth remembering…"
                    placeholderTextColor={colors.fgDisabled}
                    value={form.notes}
                    onChangeText={v => setForm(f => ({ ...f, notes: v }))}
                    onFocus={() => setFocused('notes')}
                    onBlur={() => setFocused(null)}
                    multiline
                    numberOfLines={3}
                    textAlignVertical="top"
                    textContentType="none"
                    autoComplete="off"
                  />
                </View>

                {error && (
                  <View style={styles.errorBanner}>
                    <Text style={styles.errorText}>{error}</Text>
                  </View>
                )}

                <Pressable
                  style={({ pressed }) => [
                    styles.cta,
                    (!canSave || saving) && styles.ctaDisabled,
                    pressed && canSave && !saving && styles.ctaPressed,
                  ]}
                  onPress={save}
                  disabled={!canSave || saving}
                >
                  {saving ? (
                    <ActivityIndicator color={palette.cream100} />
                  ) : (
                    <Text style={styles.ctaLabel}>
                      {isEdit ? 'Save changes' : 'Add to my gear'}
                    </Text>
                  )}
                </Pressable>
              </ScrollView>
            )}
          </KeyboardAvoidingView>
        </Animated.View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: { flex: 1, justifyContent: 'flex-end' },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(28,15,7,0.42)',
  },
  sheet: {
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    maxHeight: '92%',
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.16,
    shadowRadius: 32,
    elevation: 12,
  },
  handle: {
    width: 36,
    height: 4,
    borderRadius: 9999,
    backgroundColor: palette.cream500,
    alignSelf: 'center',
    marginTop: 14,
    marginBottom: 4,
  },

  typeContent: { paddingHorizontal: spacing[5], paddingBottom: 40 },
  sheetTitle: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 24,
    color: colors.fgPrimary,
    marginBottom: 4,
    marginTop: 8,
  },
  sheetSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: palette.espresso400,
    marginBottom: 20,
  },
  typeGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 10,
    marginTop: 20,
  },
  typeBtn: {
    width: '30%',
    flexGrow: 1,
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    borderRadius: 18,
    paddingTop: 16,
    paddingBottom: 14,
    paddingHorizontal: 8,
    alignItems: 'center',
    gap: 8,
  },
  typeBtnPressed: {
    backgroundColor: palette.cream300,
    borderColor: palette.cream500,
    transform: [{ scale: 0.95 }],
  },
  typeBtnLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: palette.espresso800,
    textAlign: 'center',
  },

  formContent: { paddingHorizontal: spacing[5], paddingBottom: 40, gap: 14 },
  backLink: { marginBottom: 4 },
  backLinkText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: palette.caramel400,
  },

  fieldWrap: { gap: 6 },
  fieldLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 12,
    color: colors.fgSecondary,
  },
  required: { color: palette.caramel500 },

  row: { flexDirection: 'row', gap: 12 },
  rowCol: { flex: 1, gap: 6 },

  input: {
    height: 50,
    paddingHorizontal: spacing[4],
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    borderRadius: 14,
    fontFamily: 'DMSans_400Regular',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  textarea: {
    height: 90,
    paddingTop: 14,
  },
  inputFocused: {
    borderColor: palette.caramel400,
    shadowColor: palette.caramel400,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 3,
  },

  errorBanner: {
    backgroundColor: palette.error100,
    borderRadius: radii.md,
    paddingVertical: 10,
    paddingHorizontal: 14,
  },
  errorText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.error500,
  },

  cta: {
    height: 54,
    borderRadius: radii.xl,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 4,
  },
  ctaDisabled: { opacity: 0.38 },
  ctaPressed: { backgroundColor: palette.espresso700 },
  ctaLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: palette.cream100,
  },
});
