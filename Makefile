##
# scop
#
# @file makefile
# @version 0.1


#       ==============================  NAME    ==============================  #
NAME            =       scop


#       ==============================  CMP     ==============================  #
CC              =       g++
CFLAGS          =       -Wall -Wextra -std=c++17 -O2 -g3 #-Werror
DFLAG           =       -MMD -MP
LDFLAGS         = -lglfw -lvulkan -ldl -lpthread -lX11 -lXxf86vm -lXrandr -lXi

#       ==============================  INCLUDE ==============================  #
DIR_INC         =       include/ ${DIR_LIBFT}
IFLAGS          =       $(DIR_INC:%=-I%)


#       ==============================  SRC     ==============================  #
SRC_DIR         =       src/
SRC_LST         =       main.cc 
SRC             = 		${addprefix ${SRC_DIR}, ${SRC_LST}}


#       ==============================  OBJ     ==============================  #
OBJ_DIR         =       .obj/
OBJ_SUBDIR      =       ${addprefix ${OBJ_DIR}, ${SRC_DIR}}


OBJ_TMP         =       ${addprefix ${OBJ_DIR}, ${SRC}}
OBJ             =       ${OBJ_TMP:.cc=.o}

DEPS            =       ${OBJ:.o=.d}




#       /*\/*\/*\/*\/*\/*\/*\/*\/*\/*\/*\/*\    RULES   /*\/*\/*\/*\/*\/*\/*\/*\/*\/*\/*\/*\    #
.PHONY: all clean fclean re

#       ==============================  BASE    ==============================  #

all: ${NAME}


clean:
		@rm -rf ${OBJ_DIR}

fclean: clean
		@rm -f ${NAME} 

re: fclean all


#       ==============================  COMPILATION     ==============================  #
${NAME}: ${OBJ_SUBDIR} ${OBJ}
		${CC} ${IFLAGS} ${CFLAGS} ${OBJ} ${LDFLAGS} -o $@



${OBJ_DIR}%.o: %.cc 
		${CC} ${IFLAGS} ${CFLAGS} ${DFLAG} -c $< -o $@



-include ${DEPS}


#       ==============================  UTILS   ==============================  #
${OBJ_SUBDIR}:
		@mkdir -p ${OBJ_SUBDIR}


# end% 