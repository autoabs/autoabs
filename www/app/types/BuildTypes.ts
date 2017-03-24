/// <reference path="../References.d.ts"/>
export const SYNC = 'build.sync';
export const UPDATE = 'build.update';
export const TRAVERSE = 'build.traverse';
export const REMOVE = 'build.remove';
export const CHANGE = 'build.change';

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
	repo_state?: string;
	arch?: string;
	pkg_ids?: string[];
	pkg_build_id?: string;
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
