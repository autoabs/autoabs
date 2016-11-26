/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Mui from 'material-ui';
import Styles from '../Styles';
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
		color: Styles.colors.color,
		margin: '1px 1px 0 0',
	} as React.CSSProperties,
	content: {
		padding: '10px 0 10px 10px',
	} as React.CSSProperties,
	name: {
		fontSize: '20px',
	} as React.CSSProperties,
	version: {
		fontSize: '14px',
		margin: '6px 0 0 7px',
		color: Styles.colors.fadeColor,
	} as React.CSSProperties,
	repo: {
		fontSize: '12px',
		marginTop: '7px',
		color: Styles.colors.fadeColor,
	} as React.CSSProperties,
	actions: {
		padding: '5px 0',
		justifyContent: 'center',
	} as React.CSSProperties,
	pause: {
		display: 'none',
		color: styles.colors.blue500,
	} as React.CSSProperties,
	resume: {
		display: 'none',
		color: styles.colors.green500,
	} as React.CSSProperties,
	retry: {
		color: styles.colors.pink400,
	} as React.CSSProperties,
	remove: {
		color: styles.colors.red500,
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

		return <Mui.Card style={css.card}>
			<Mui.CardText>
				<div className="layout horizontal">
					<div style={css.content} className="card-content flex">
						<div className="layout vertical">
							<div className="layout horizontal">
								<div style={css.name}>{build.name}</div>
								<div style={css.version}>
									{build.version}-{build.release}
								</div>
							</div>
							<div style={css.repo}>{build.repo} - {build.arch}</div>
						</div>
					</div>
					<div>
						<Mui.IconButton style={css.launch}>
							<i className="material-icons">flip_to_front</i>
						</Mui.IconButton>
					</div>
				</div>
			</Mui.CardText>
			<Mui.CardActions style={css.actions}>
				<div className="layout horizontal">
					<Mui.FlatButton style={css.pause} label="Pause"/>
					<Mui.FlatButton style={css.resume} label="Resume"/>
					<Mui.FlatButton style={css.retry} label="Retry"/>
					<Mui.FlatButton style={css.remove} label="Remove"/>
				</div>
			</Mui.CardActions>
		</Mui.Card>;
	}
}
