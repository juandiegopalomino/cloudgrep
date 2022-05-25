import { CSSProperties } from 'react';

export interface SearchProps {
	onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	width?: string;
	height?: string;
	rest?: CSSProperties;
}
