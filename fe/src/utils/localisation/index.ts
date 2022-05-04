import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

i18n.use(initReactI18next).init({
	resources: {
		en: {
			translation: {
				APP_NAME: 'Cloud grep',

				// GENERAL
				ADD: 'Add',
				SAVE: 'Save',
				CANCEL: 'Cancel',
				DELETE: 'Delete',
				Edit: 'Edit',
				TRY_AGAIN: 'Try Again',
				SOMETHING_WENT_WRONG: 'Something went wrong',
				LOGIN: 'Login',
				LOGOUT: 'Logout',
				HAS_NO_ACC: 'Do not have an account?',
				REGISTER: 'Register',
				EXISTING_ACCOUNT: 'Already have an account?',
				NOT_FOUND: 'Not Found',

				// Application specific
				TAGS: 'Tags',
			},
		},
	},
	lng: 'en',
	debug: false,
	keySeparator: false,
	interpolation: {
		escapeValue: false,
		formatSeparator: ',',
	},
	react: {
		useSuspense: false,
	},
});

export default i18n;
