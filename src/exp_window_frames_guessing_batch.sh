NCPU=4
B_MAX=32
START=`date +%s`

for G in 5 6 7 8 9; do
	for N in 5 8 10 12 15 16 17 18 20; do
		for R in 4 8 12 16 20; do
			if [[ ! -e './G0'$G'N'$N'R'$R ]]; then
				mkdir './G0'$G'N'$N'R'$R
			fi
			if [[ ! -f './G0'$G'N'$N'R'$R'/G0'$G'N'$N'R'$R'_parameters.json' ]]; then
				go run ans_parameters_generator.go -G 0.$G -N $N -R $R -prefix 'G0'$G'N'$N'R'$R'/G0'$G'N'$N'R'$R
			fi
			if [[ ! -f './G0'$G'N'$N'R'$R'/G0'$G'N'$N'R'$R'_config.json' ]]; then
				go run ans_initialization.go -prefix 'G0'$G'N'$N'R'$R/'G0'$G'N'$N'R'$R
			fi
			for S in {4..16}; do
				for (( i = 2; i <= $B_MAX; i+=$NCPU )); do
					for (( B = $i; B < ($i + $NCPU); B++ )); do
						if [[ $B -gt $B_MAX ]]; then
							break
						fi
						if [[ ! -f './G0'$G'N'$N'R'$R'/G0'$G'N'$N'R'$R'_S'$S'_B'$B'_guessing_window_frame.json' ]]; then
							printf "G 0.%d\t N %d\t R %d\t S %02d\t B %02d\n" $G $N $R $S $B
							go run exp_window_frames_guessing.go -prefix 'G0'$G'N'$N'R'$R'/G0'$G'N'$N'R'$R -s $S -b $B &
						fi
					done
					wait
				done
			done
		done
	done
done
END=`date +%s`
RUNTIME=$(( $END - $START ))
echo Finished in $RUNTIME seconds
