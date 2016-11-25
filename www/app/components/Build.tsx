/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildTypes from '../types/BuildTypes';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	build: BuildTypes.Build;
}

const css = {
	card: {
		flex: '1 0 auto',
		minWidth: '300px',
		maxWidth: '450px',
		margin: '5px',
	} as React.CSSProperties,
	launch: {
		color: '#7f7f7f',
		margin: '1px 1px 0 0',
	} as React.CSSProperties,
	content: {
		paddingRight: '0',
	} as React.CSSProperties,
	name: {
		fontSize: '20px',
	} as React.CSSProperties,
	version: {
		fontSize: '14px',
		margin: '6px 0 0 7px',
		color: '#919191',
	} as React.CSSProperties,
	repo: {
		fontSize: '12px',
		marginTop: '7px',
		color: '#919191',
	} as React.CSSProperties,
	actions: {
		justifyContent: 'center',
	} as React.CSSProperties,
	logs: {
		color: '#607d8b',
	} as React.CSSProperties,
	pause: {
		display: 'none',
		color: '#03a9f4',
	} as React.CSSProperties,
	resume: {
		display: 'none',
		color: '#4caf50',
	} as React.CSSProperties,
	retry: {
		color: '#ab47bc',
	} as React.CSSProperties,
	remove: {
		color: '#f44336',
	} as React.CSSProperties,
};

export default class Build extends React.Component<Props, null> {
	render(): JSX.Element {
		let build = this.props.build;

		let start = '-';
		if (build.start !== "0001-01-01T00:00:00Z") {
			start = MiscUtils.formatDate(new Date(build.start));
		}
		let stop = '-';
		if (build.stop !== "0001-01-01T00:00:00Z") {
			stop = MiscUtils.formatDate(new Date(build.stop));
		}

		return <paper-card style={css.card} alt={build.name}>
			<div className="layout horizontal">
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div className="layout horizontal">
							<div style={css.name}>{build.name}</div>
							<div style={css.version}>{build.version}-{build.release}</div>
						</div>
						<div style={css.repo}>{build.repo} - {build.arch}</div>
					</div>
				</div>
				<paper-icon-button style={css.launch} icon="launch"/>
			</div>
			<div style={css.actions} className="card-actions layout horizontal">
				<paper-button style={css.pause}>Pause</paper-button>
				<paper-button style={css.resume}>Resume</paper-button>
				<paper-button style={css.retry}>Retry</paper-button>
				<paper-button style={css.remove}>Remove</paper-button>
			</div>
		</paper-card>;
	}
}
