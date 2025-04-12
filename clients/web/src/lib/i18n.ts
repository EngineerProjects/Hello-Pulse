export type Locale = "fr" | "en";

export const DEFAULT_LOCALE: Locale = "fr";

export const locales: Record<Locale, Record<string, string>> = {
  fr: {
    "auth.login": "Se connecter",
    "auth.register": "Créer un compte",
    "auth.or": "OU",
    "auth.firstName": "Prénom",
    "auth.lastName": "Nom",
    "auth.email": "Mail",
    "auth.password": "Mot de passe",
    "auth.confirmPassword": "Confirmer le mot de passe",
    "auth.submit.login": "Connexion",
    "auth.submit.register": "S'enregistrer",
    "auth.forgotPassword": "Mot de passe oublié ?",
    "auth.alreadyHaveAccount": "Vous avez déjà un compte ?",
    "auth.termsOfUse": "Conditions d'utilisation",
    "auth.privacyPolicy": "Politique de confidentialité",
    "auth.title": "À vos marque, prêt, feu, go ! Brainstormer",
    "auth.description": "Entrez dans votre espace créatif et donnez vie à vos idées avec Hello Pulse. Connectez-vous pour rejoindre la prochaine session de brainstorming en direct, où l'inspiration est sans limites et les idées fusent en un éclair !",
  },
  en: {
    "auth.login": "Log In",
    "auth.register": "Create an Account",
    "auth.or": "OR",
    "auth.firstName": "First Name",
    "auth.lastName": "Last Name",
    "auth.email": "Email Address",
    "auth.password": "Password",
    "auth.confirmPassword": "Confirm Password",
    "auth.submit.login": "Log In",
    "auth.submit.register": "Register",
    "auth.forgotPassword": "Forgot Password?",
    "auth.alreadyHaveAccount": "Already have an account?",
    "auth.termsOfUse": "Terms of Use",
    "auth.privacyPolicy": "Privacy Policy",
    "auth.title": "On your marks, get set, go! Brainstorm!",
    "auth.description": "Step into your creative space and bring your ideas to life with Hello Pulse. Log in to join the next live brainstorming session, where inspiration knows no limits and ideas flow in a flash!",
  },
};

export function getTranslation(key: string, locale: Locale = DEFAULT_LOCALE): string {
  return locales[locale][key] || key;
}