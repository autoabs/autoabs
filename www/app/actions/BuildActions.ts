/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import * as BuildTypes from '../types/BuildTypes';

export function sync(): Promise<string> {
	Dispatcher.dispatch({
		type: BuildTypes.LOADING,
	});

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.get('/build')
			.query({'index': 0})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				Dispatcher.dispatch({
					type: BuildTypes.LOADED,
				});

				if (err) {
					Alert.error('Failed to sync builds');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: BuildTypes.SYNC,
					data: {
						builds: res.body.builds,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function archive(id: string): Promise<string> {
	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.put('/build/' + id + '/archive')
			.set('Accept', 'application/json')
			.end((err: any): void => {
				Dispatcher.dispatch({
					type: BuildTypes.LOADED,
				});

				if (err) {
					Alert.error('Failed to archive build');
					reject(err);
					return;
				}

				sync().then(resolve, resolve);
			});
	});
}

export function rebuild(id: string): Promise<string> {
	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.put('/build/' + id + '/rebuild')
			.set('Accept', 'application/json')
			.end((err: any): void => {
				Dispatcher.dispatch({
					type: BuildTypes.LOADED,
				});

				if (err) {
					Alert.error('Failed to rebuild build');
					reject(err);
					return;
				}

				sync().then(resolve, resolve);
			});
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: BuildTypes.REMOVE,
		data: {
			id: id,
		},
	});

	Alert.info('Build successfully removed');
}
