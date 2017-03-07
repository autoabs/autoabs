/// <reference path="../References.d.ts"/>
export const SYNC = 'build.sync';
export const UPDATE = 'build.update';
export const REMOVE = 'build.remove';

export interface Build {
	id: string;
	name?: string;
	builder?: string;
	start?: string;
	stop?: string;
	state?: string;
	version?: string;
	release?: string;
	repo?: string;
	arch?: string;
}

export type Builds = Build[];

export interface BuildDispatch {
	type: string;
	data?: {
		id?: string;
		content?: string;
		build?: Build;
		builds?: Build[];
		index?: number;
		count?: number;
	};
}
