import {
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
  ActivityIndicator,
} from 'react-native';
import { useRef, useState } from 'react';
import { router } from 'expo-router';
import Svg, { Defs, RadialGradient, Stop, Path } from 'react-native-svg';
import { colors, spacing, radii } from '@/src/theme';
import { apiFetch } from '@/src/config/api';
import { saveTokens } from '@/src/lib/auth';

// ─── Types ────────────────────────────────────────────────────────────────────

type Mode = 'login' | 'register';

interface Form {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  confirm: string;
}

interface AuthResponse {
  token: string;
  refresh_token: string;
}

// ─── Bean mark SVG ────────────────────────────────────────────────────────────

function BeanMark({ size = 48 }: { size?: number }) {
  const h = size * 1.22;
  return (
    <Svg width={size} height={h} viewBox="0 0 36 44">
      <Defs>
        <RadialGradient id="bmbg" cx="38%" cy="30%" r="65%" gradientUnits="userSpaceOnUse" fx="38%" fy="30%">
          <Stop offset="0%" stopColor="#3A1E10" />
          <Stop offset="100%" stopColor="#0F0603" />
        </RadialGradient>
      </Defs>
      {/* Steam wisps */}
      <Path d="M9 9 C7.5 6 10 3.5 8.5 1"   stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      <Path d="M18 7 C16.5 4 19 1.5 17.5 -1" stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      <Path d="M27 9 C25.5 6 28 3.5 26.5 1" stroke="#D4872A" strokeWidth="1.8" strokeLinecap="round" />
      {/* Bean body */}
      <Path
        d="M13 12 C6 14 2 21 2 28 C2 36 7 44 14 44 C17 44.5 19 44.5 22 44 C29 43 34 36 34 28 C34 20 29 13 22 12 C19 11 16 11 13 12 Z"
        fill="url(#bmbg)"
      />
      {/* Bean crease */}
      <Path d="M21 13 C15 23 23 31 15 40 C13 43 16 44 18 44" stroke="#FAF3E8" strokeWidth="1.5" strokeLinecap="round" opacity="0.5" />
    </Svg>
  );
}

// ─── Screen ───────────────────────────────────────────────────────────────────

export default function LoginScreen() {
  const [mode, setMode] = useState<Mode>('login');
  const [form, setForm] = useState<Form>({ firstName: '', lastName: '', email: '', password: '', confirm: '' });
  const [focusedField, setFocusedField] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const lastNameRef  = useRef<TextInput>(null);
  const emailRef     = useRef<TextInput>(null);
  const passwordRef  = useRef<TextInput>(null);
  const confirmRef   = useRef<TextInput>(null);

  function update(key: keyof Form, value: string) {
    setForm(f => ({ ...f, [key]: value }));
    setError(null);
  }

  function switchMode(next: Mode) {
    setMode(next);
    setError(null);
  }

  const canSubmit =
    mode === 'login'
      ? form.email.trim().length > 0 && form.password.length > 0
      : form.firstName.trim().length > 0 &&
        form.lastName.trim().length > 0 &&
        form.email.trim().length > 0 &&
        form.password.length >= 8 &&
        form.password === form.confirm;

  async function handleSubmit() {
    if (loading || !canSubmit) return;

    if (mode === 'register') {
      if (form.password.length < 8) { setError('Password must be at least 8 characters.'); return; }
      if (form.password !== form.confirm) { setError("Passwords don't match."); return; }
    }

    setLoading(true);
    setError(null);

    try {
      let res: AuthResponse;
      if (mode === 'login') {
        res = await apiFetch<AuthResponse>('/auth/login', {
          method: 'POST',
          body: { email: form.email.trim(), password: form.password },
        });
      } else {
        res = await apiFetch<AuthResponse>('/auth/register', {
          method: 'POST',
          body: {
            first_name: form.firstName.trim(),
            last_name: form.lastName.trim(),
            email: form.email.trim(),
            password: form.password,
          },
        });
      }
      await saveTokens(res.token, res.refresh_token);
      router.replace('/(tabs)');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went sideways. Even my grinder jams sometimes.');
    } finally {
      setLoading(false);
    }
  }

  function inputStyle(field: string) {
    return [styles.input, focusedField === field && styles.inputFocused];
  }

  return (
    <KeyboardAvoidingView
      style={styles.screen}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView
        contentContainerStyle={styles.scroll}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        {/* Brand hero */}
        <View style={styles.hero}>
          <BeanMark size={48} />
          <Text style={styles.brandTitle}>Mr. Bean</Text>
          <Text style={styles.brandSubtitle}>Your espresso, perfected.</Text>
        </View>

        {/* Segment control */}
        <View style={styles.segmentTrack}>
          {(['login', 'register'] as const).map(m => (
            <Pressable
              key={m}
              onPress={() => switchMode(m)}
              style={[styles.segmentItem, mode === m && styles.segmentActive]}
            >
              <Text style={[styles.segmentLabel, mode === m && styles.segmentLabelActive]}>
                {m === 'login' ? 'Sign in' : 'Create account'}
              </Text>
            </Pressable>
          ))}
        </View>

        {/* Form */}
        <View style={styles.form}>
          {mode === 'register' && (
            <View style={styles.nameRow}>
              <View style={styles.nameCol}>
                <Text style={styles.fieldLabel}>First name <Text style={styles.required}>*</Text></Text>
                <TextInput
                  style={inputStyle('firstName')}
                  placeholder="Ada"
                  placeholderTextColor={colors.fgDisabled}
                  value={form.firstName}
                  onChangeText={v => update('firstName', v)}
                  autoCapitalize="words"
                  returnKeyType="next"
                  onSubmitEditing={() => lastNameRef.current?.focus()}
                  onFocus={() => setFocusedField('firstName')}
                  onBlur={() => setFocusedField(null)}
                />
              </View>
              <View style={styles.nameCol}>
                <Text style={styles.fieldLabel}>Last name <Text style={styles.required}>*</Text></Text>
                <TextInput
                  ref={lastNameRef}
                  style={inputStyle('lastName')}
                  placeholder="Lovelace"
                  placeholderTextColor={colors.fgDisabled}
                  value={form.lastName}
                  onChangeText={v => update('lastName', v)}
                  autoCapitalize="words"
                  returnKeyType="next"
                  onSubmitEditing={() => emailRef.current?.focus()}
                  onFocus={() => setFocusedField('lastName')}
                  onBlur={() => setFocusedField(null)}
                />
              </View>
            </View>
          )}

          <View>
            <Text style={styles.fieldLabel}>Email <Text style={styles.required}>*</Text></Text>
            <TextInput
              ref={emailRef}
              style={inputStyle('email')}
              placeholder="you@example.com"
              placeholderTextColor={colors.fgDisabled}
              value={form.email}
              onChangeText={v => update('email', v)}
              keyboardType="email-address"
              autoCapitalize="none"
              autoCorrect={false}
              returnKeyType="next"
              onSubmitEditing={() => passwordRef.current?.focus()}
              onFocus={() => setFocusedField('email')}
              onBlur={() => setFocusedField(null)}
            />
          </View>

          <View>
            <Text style={styles.fieldLabel}>Password <Text style={styles.required}>*</Text></Text>
            <TextInput
              ref={passwordRef}
              style={inputStyle('password')}
              placeholder={mode === 'register' ? 'Min. 8 characters' : 'Your password'}
              placeholderTextColor={colors.fgDisabled}
              value={form.password}
              onChangeText={v => update('password', v)}
              secureTextEntry
              returnKeyType={mode === 'register' ? 'next' : 'done'}
              onSubmitEditing={() => mode === 'register' ? confirmRef.current?.focus() : handleSubmit()}
              onFocus={() => setFocusedField('password')}
              onBlur={() => setFocusedField(null)}
            />
          </View>

          {mode === 'register' && (
            <View>
              <Text style={styles.fieldLabel}>Confirm password <Text style={styles.required}>*</Text></Text>
              <TextInput
                ref={confirmRef}
                style={inputStyle('confirm')}
                placeholder="Same again"
                placeholderTextColor={colors.fgDisabled}
                value={form.confirm}
                onChangeText={v => update('confirm', v)}
                secureTextEntry
                returnKeyType="done"
                onSubmitEditing={handleSubmit}
                onFocus={() => setFocusedField('confirm')}
                onBlur={() => setFocusedField(null)}
              />
            </View>
          )}

          {error && (
            <View style={styles.errorBanner}>
              <Text style={styles.errorText}>{error}</Text>
            </View>
          )}

          <Pressable
            style={({ pressed }) => [
              styles.cta,
              (!canSubmit || loading) && styles.ctaDisabled,
              pressed && canSubmit && !loading && styles.ctaPressed,
            ]}
            onPress={handleSubmit}
            disabled={!canSubmit || loading}
          >
            {loading ? (
              <ActivityIndicator color={colors.fgInverse} />
            ) : (
              <Text style={styles.ctaLabel}>
                {mode === 'login' ? 'Sign in' : 'Create account'}
              </Text>
            )}
          </Pressable>

          {mode === 'login' && (
            <Pressable style={styles.forgotWrap}>
              <Text style={styles.forgotLabel}>Forgot password?</Text>
            </Pressable>
          )}
        </View>

        <View style={styles.bottomPad} />
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bgApp,
  },
  scroll: {
    flexGrow: 1,
  },

  // Hero
  hero: {
    paddingTop: spacing[10],
    paddingHorizontal: spacing[8],
    paddingBottom: spacing[8],
    alignItems: 'center',
  },
  brandTitle: {
    fontFamily: 'PlayfairDisplay_900Black',
    fontSize: 36,
    lineHeight: 40,
    letterSpacing: -0.8,
    color: colors.fgPrimary,
    marginTop: spacing[4] + 2,
  },
  brandSubtitle: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    lineHeight: 21,
    color: colors.fgSecondary,
    marginTop: spacing[2],
  },

  // Segment
  segmentTrack: {
    flexDirection: 'row',
    backgroundColor: colors.bgSubtle,
    borderRadius: 14,
    padding: 4,
    marginHorizontal: spacing[6],
    marginBottom: 28,
  },
  segmentItem: {
    flex: 1,
    height: 36,
    borderRadius: 10,
    alignItems: 'center',
    justifyContent: 'center',
  },
  segmentActive: {
    backgroundColor: colors.bgApp,
    shadowColor: '#1C0F07',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  segmentLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: colors.fgSecondary,
  },
  segmentLabelActive: {
    color: colors.fgPrimary,
  },

  // Form
  form: {
    paddingHorizontal: spacing[6],
    gap: 14,
  },
  nameRow: {
    flexDirection: 'row',
    gap: 12,
  },
  nameCol: {
    flex: 1,
  },
  fieldLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: colors.fgSecondary,
    marginBottom: 6,
  },
  required: {
    color: colors.fgAccent,
  },
  input: {
    height: 50,
    paddingHorizontal: spacing[4],
    backgroundColor: colors.bgCard,
    borderWidth: 1.5,
    borderColor: colors.borderSubtle,
    borderRadius: 14,
    fontFamily: 'DMSans_400Regular',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  inputFocused: {
    borderColor: colors.borderFocus,
    shadowColor: colors.borderFocus,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 3,
    elevation: 0,
  },

  // Error
  errorBanner: {
    backgroundColor: '#FDECEA',
    borderWidth: 1,
    borderColor: '#F5C6C2',
    borderRadius: radii.md,
    paddingVertical: 10,
    paddingHorizontal: 14,
  },
  errorText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: '#C0392B',
  },

  // CTA
  cta: {
    height: 54,
    borderRadius: radii.xl,
    backgroundColor: colors.interactivePrimary,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 4,
  },
  ctaDisabled: {
    opacity: 0.38,
  },
  ctaPressed: {
    transform: [{ scale: 0.97 }],
    backgroundColor: colors.interactivePrimaryHover,
  },
  ctaLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: colors.fgInverse,
  },

  // Forgot
  forgotWrap: {
    alignItems: 'center',
    paddingVertical: 2,
  },
  forgotLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: colors.fgAccent,
  },

  bottomPad: {
    height: spacing[12],
  },
});
