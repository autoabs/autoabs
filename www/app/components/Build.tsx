/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Card from 'material-ui/Card';
import CardText from 'material-ui/Card';
import CardActions from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
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
		padding: '10px',
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

		return <Card style={css.card}>
			<CardText className="layout horizontal">
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div className="layout horizontal">
							<div style={css.name}>{build.name}</div>
							<div style={css.version}>{build.version}-{build.release}</div>
						</div>
						<div style={css.repo}>{build.repo} - {build.arch}</div>
					</div>
				</div>
			</CardText>
			<CardActions style={css.actions} className="layout horizontal">
				<FlatButton style={css.pause} label="Pause"/>
				<FlatButton style={css.resume} label="Resume"/>
				<FlatButton style={css.retry} label="Retry"/>
				<FlatButton style={css.remove} label="Remove"/>
			</CardActions>
		</Card>;
	}
}
