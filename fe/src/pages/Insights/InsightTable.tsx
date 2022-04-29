import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { Property } from 'models/Resource';
import React, { FC } from 'react';
import { useAppSelector } from 'store/hooks';

const InsightTable: FC = () => {
	const { tags } = useAppSelector(state => state.tags);
	const { resources } = useAppSelector(state => state.resources);

	return (
		<Box
			sx={{
				width: '80%',
				height: '100%',
				'&:hover': {
					backgroundColor: 'primary.main',
					opacity: [0.9, 0.8, 0.7],
				},
			}}>
			<TableContainer component={Paper}>
				<Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
					<TableHead>
						<TableRow>
							<TableCell>Type </TableCell>
							<TableCell align="right">Id</TableCell>
							<TableCell align="right">Region</TableCell>
							<TableCell align="right">Name</TableCell>
							<TableCell align="right">Value</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{resources.map(row =>
							row.Properties?.map((resource: Property) => (
								<TableRow key={row.Id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
									<TableCell component="th" scope="row">
										{row.Type}
									</TableCell>
									<TableCell align="right">{row.Id}</TableCell>
									<TableCell align="right">{row.Region}</TableCell>
									<TableCell align="right">{resource.Name}</TableCell>
									<TableCell align="right">{resource.Value}</TableCell>
								</TableRow>
							))
						)}
					</TableBody>
				</Table>
			</TableContainer>
		</Box>
	);
};

export default InsightTable;