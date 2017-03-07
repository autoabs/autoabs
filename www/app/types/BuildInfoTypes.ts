/// <reference path="../References.d.ts"/>
export const OPEN = 'build_info.open';
export const CLOSE = 'build_info.close';
export const UPDATE = 'build_info.update';

export interface BuildInfoDispatch {
	type: string;
	data?: {
		id?: string;
		output?: string[];
	};
}
